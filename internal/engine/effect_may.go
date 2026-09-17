package engine

import (
	"fmt"
	"reflect"
)

// A "you may" effect is optional: it offers the controller the choice to resolve
// its inner effect or to decline it entirely. It models KeyForge's "You may <do
// X>", where passing is always allowed even when a legal target exists — the
// distinction that keeps Chuff Ape's "you may destroy another friendly creature"
// from ever being forced.
type May struct {
	Do Effect
}

// Text renders the effect, e.g. "you may destroy another friendly creature -> fully
// heal Chuff Ape".
func (e May) Text() string {
	return "you may " + e.Do.Text()
}

// A declinableEffect is an Effect whose whole decision is "which card — or none":
// a single target chosen from a pool. May prefers it over its own Yes/No question,
// so a player told "you may destroy another friendly creature" clicks the creature
// they mean (or passes) instead of answering twice. Like GatingEffect the method
// is unexported, so only engine effects can offer the shortcut, and it reports
// whether anything happened so a gate wrapping it still works.
type declinableEffect interface {
	Effect
	// declinable reports whether this effect's decision really is a single card
	// choice. "You may destroy each Mars creature" is not — there is nothing to
	// click — so it keeps the Yes/No.
	declinable() bool
	// resolveOptional resolves the effect as its own optional choice, returning
	// whether anything happened.
	resolveOptional(ctx *EffectContext) bool
}

// A vacuousEffect can tell, before it is offered, that it would do nothing at
// all. May uses it to skip a question with only one honest answer: "you may
// destroy each Mars creature" is not a decision when no Mars creature is in play.
type vacuousEffect interface {
	Effect
	vacuous(ctx *EffectContext) bool
}

// Resolve offers the inner effect. When that effect is itself one declinable card
// choice, the player makes that choice directly; when there is no card to choose
// but the card asking is on the board, they opt in by clicking that card;
// otherwise they answer Yes/No and the effect resolves on Yes. An effect with
// nothing to act on is not offered.
func (e May) Resolve(ctx *EffectContext) {
	if d, ok := e.Do.(declinableEffect); ok && d.declinable() {
		d.resolveOptional(ctx)
		return
	}
	if v, ok := e.Do.(vacuousEffect); ok && v.vacuous(ctx) {
		return
	}
	prompt := capitalizeFirst(e.Text())
	if e.offeredOnSource(ctx) {
		if _, ok := ctx.ChooseCardOptional(prompt, []LocalID{ctx.Source}); !ok {
			return
		}
		e.Do.Resolve(ctx)
		return
	}
	if ctx.ChooseOption(prompt, []string{"Yes", "No"}) == 0 {
		e.Do.Resolve(ctx)
	}
}

// offeredOnSource reports that the offer can be opted into by clicking the card
// that makes it, so an optional ability with no card of its own to choose still
// reads as a click rather than a Yes/No (Gambling Den: click it to be asked for a
// house, or Done to pass). It needs the card to be on the board — a Tactic has
// already left play when its Play ability resolves (The Common Cold) and a card
// triggering from the discard was never there (Relentless Creeper). A combatant
// mid-fight is excluded too: the player just clicked it to fight, so clicking it
// again would not read as opting in (Angry Mob).
func (e May) offeredOnSource(ctx *EffectContext) bool {
	return ctx.Resolver.InPlay(ctx.Source) && !ctx.Resolver.CurrentlyFighting(ctx.Source)
}

// validate descends into the wrapped effect, and holds the shape rule that keeps
// declinability from being forgotten: an effect whose decision is a chosen card
// must offer that click, not a Yes/No. Checking it here means a card that would
// ask twice fails at init (ADR 0010) rather than only being noticed as an odd
// prompt in play, and a new single-card effect cannot quietly regrow the
// inconsistency the way Exalt once did against Destroy.
func (e May) validate() error {
	if err := validateEffect(e.Do); err != nil {
		return err
	}
	if _, ok := e.Do.(declinableEffect); ok {
		return nil
	}
	if chosenTargetField(e.Do) {
		return fmt.Errorf(
			"May: %T decides by choosing a card but is not declinable; "+
				"implement declinable/resolveOptional so the May is answered by the click",
			e.Do,
		)
	}
	return nil
}

// chosenTargetField reports whether an effect's own fields hold a Target the
// controller chooses — the shape that makes the effect a single-card decision. It
// reads the struct rather than asking the effect, so the answer does not depend on
// the effect remembering to say so. A configured exception (DealDamage with a
// Spread) still implements declinableEffect and answers false there.
func chosenTargetField(e Effect) bool {
	v := reflect.ValueOf(e)
	if v.Kind() != reflect.Struct {
		return false
	}
	for i := range v.NumField() {
		f := v.Field(i)
		if t, ok := f.Interface().(Target); ok && t.isChosen() {
			return true
		}
	}
	return false
}
