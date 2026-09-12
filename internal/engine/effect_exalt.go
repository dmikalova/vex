package engine

import (
	"errors"
	"fmt"
)

// To exalt a creature is to place 1 Æmber from the common supply onto a chosen
// friendly or enemy creature. The Æmber sits on the creature, belonging to no
// pool, until it leaves play, then goes to the owner's opponent's pool. Exalting
// N times places N Æmber.
type Exalt struct {
	Target Target
	Amount int
}

// validate requires an explicit target.
func (e Exalt) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Exalt")
	}
	return nil
}

// Text renders the effect, e.g. "exalt an enemy creature 2 times". A single
// exalt drops the count so it reads naturally.
func (e Exalt) Text() string {
	if e.Amount == 1 {
		return "exalt " + e.Target.Text()
	}
	return fmt.Sprintf("exalt %s %d times", e.Target.Text(), e.Amount)
}

// verb and targetText let a single exalt fold with a neighbouring combinable on
// the same target in a Sequence — "ready and exalt each neighboring creature".
func (e Exalt) verb() string       { return "exalt" }
func (e Exalt) targetText() string { return e.Target.Text() }

// foldable folds only a single exalt; "exalt X 2 times" keeps its own shape.
func (e Exalt) foldable() bool { return e.Amount == 1 }

// Resolve chooses a creature (through the Target) and places Amount Æmber on it.
func (e Exalt) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate exalts each selected creature and reports whether any was exalted,
// so an exalt can be the first half of a Then — Saurus Rex only searches its deck
// when the optional exalt happens.
func (e Exalt) resolveGate(ctx *EffectContext) bool {
	ids := e.Target.Select(ctx)
	for _, id := range ids {
		ctx.Resolver.AddAmberOn(id, e.Amount)
		ctx.Resolver.Record(AemberExalted{Creature: id, Amount: e.Amount})
	}
	return len(ids) > 0
}

// declinable reports that a self-exalt is one clickable card — the source — so
// "you may exalt Senator Shrix" is answered by clicking Shrix or passing, not by a
// separate Yes/No.
func (e Exalt) declinable() bool { return e.Target.Kind == TargetThisCreature }

// resolveOptional is Resolve under a May: the source is offered as a clickable
// card, so the exalt is confirmed by clicking it or declined with Done.
func (e Exalt) resolveOptional(ctx *EffectContext) bool {
	id, ok := ctx.ChooseCardOptional("Exalt "+e.Target.Text(), e.Target.Select(ctx))
	if !ok {
		return false
	}
	ctx.Resolver.AddAmberOn(id, e.Amount)
	ctx.Resolver.Record(AemberExalted{Creature: id, Amount: e.Amount})
	return true
}

// ExaltToRepeat resolves Do once, then lets the controller exalt a creature to
// repeat it a single time — the exalt pays for one repeat. "Repeat the preceding
// effect" repeats only the effect before the exalt clause, so the offer is made
// once and does not chain. It models "<do>. You may exalt <a creature> to repeat
// the preceding effect." (Phalanx Strike, Tribute).
type ExaltToRepeat struct {
	// Do is the preceding effect each exalt repeats.
	Do Effect
	// Exalt names the creature exalted to pay for a repeat.
	Exalt Target
}

// validate requires an effect to repeat and an explicit exalt target.
func (e ExaltToRepeat) validate() error {
	if e.Do == nil {
		return errors.New("ExaltToRepeat needs an effect to Do")
	}
	if !e.Exalt.valid() {
		return errUnsetTarget("ExaltToRepeat")
	}
	return validateEffect(e.Do)
}

// Text renders Do as its own sentence, then the optional exalt-to-repeat gate.
func (e ExaltToRepeat) Text() string {
	return punctuate(e.Do.Text()) +
		" You may exalt " + e.Exalt.Text() + " to repeat the preceding effect"
}

// Resolve runs Do once, then offers a single exalt-to-repeat (a declinable
// choice). Declining exalts nothing and ends there; accepting exalts the chosen
// creature and resolves Do one more time.
func (e ExaltToRepeat) Resolve(ctx *EffectContext) {
	e.Do.Resolve(ctx)
	ids := e.exaltChoice(ctx)
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		ctx.Resolver.AddAmberOn(id, 1)
		ctx.Resolver.Record(AemberExalted{Creature: id, Amount: 1})
	}
	e.Do.Resolve(ctx)
}

// exaltChoice offers the exalt that pays for another repeat, or none to stop. A
// chosen target is its own declinable prompt (pick a creature or decline); a
// back-reference like the chosen creature has nothing to pick, so it is offered
// as a Yes/No confirm on the creature Do just acted on.
func (e ExaltToRepeat) exaltChoice(ctx *EffectContext) []LocalID {
	if e.Exalt.isChosen() {
		return e.Exalt.SelectOptional(ctx)
	}
	ids := e.Exalt.Select(ctx)
	if len(ids) == 0 {
		return nil
	}
	prompt := "Exalt " + e.Exalt.Text() + " to repeat the preceding effect?"
	if ctx.ChooseOption(prompt, []string{"Yes", "No"}) != 0 {
		return nil
	}
	return ids
}
