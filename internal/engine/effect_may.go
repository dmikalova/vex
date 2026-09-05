package engine

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
// choice, the player makes that choice directly; otherwise they answer Yes/No and
// the effect resolves on Yes. An effect with nothing to act on is not offered.
func (e May) Resolve(ctx *EffectContext) {
	if d, ok := e.Do.(declinableEffect); ok && d.declinable() {
		d.resolveOptional(ctx)
		return
	}
	if v, ok := e.Do.(vacuousEffect); ok && v.vacuous(ctx) {
		return
	}
	if ctx.ChooseOption(capitalizeFirst(e.Text()), []string{"Yes", "No"}) == 0 {
		e.Do.Resolve(ctx)
	}
}

// validate descends into the wrapped effect.
func (e May) validate() error { return validateEffect(e.Do) }
