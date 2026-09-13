package engine

import "fmt"

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
