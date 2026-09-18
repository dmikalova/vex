package engine

import (
	"fmt"
	"slices"
)

// To exalt a creature is to place 1 Æmber from the common supply onto a chosen
// friendly or enemy creature. The Æmber sits on the creature, belonging to no
// pool, until it leaves play, then goes to the owner's opponent's pool. Exalting
// N times places N Æmber. The two count axes are distinct: Amount scales how much
// Æmber one exalt places, while Times repeats the exalt, choosing a fresh Target
// each time (Bawretchadontius exalts two enemy creatures).
type Exalt struct {
	Target Target
	Amount int
	// Times repeats the exalt, choosing a fresh Target each time.
	Times Count
	// Distinct bars a creature an earlier repetition already exalted from being
	// picked again, so a Times that repeats N times spreads the exalts across N
	// different creatures (Bawretchadontius exalts two distinct enemy creatures).
	Distinct bool
}

// validate requires an explicit target, and a Times for a Distinct spread.
func (e Exalt) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Exalt")
	}
	if e.Distinct && e.Times == nil {
		return fmt.Errorf("Exalt: Distinct is meaningless without a Times to repeat")
	}
	return nil
}

// Text renders the effect, e.g. "exalt an enemy creature", "exalt an enemy
// creature 2 times" (repeating a single target's exalt), or "exalt 2 enemy
// creatures" when Times spreads the exalt across that many distinct creatures.
func (e Exalt) Text() string { return "exalt " + e.targetText() }

// verb and targetText let a single exalt fold with a neighbouring combinable on
// the same target in a Sequence — "exalt a friendly creature and 2 enemy
// creatures". targetText names what is exalted: the plain target, a repeat count
// ("an enemy creature 2 times"), or a distinct spread ("2 enemy creatures").
func (e Exalt) verb() string { return "exalt" }
func (e Exalt) targetText() string {
	if n := fixedValue(e.Times); n > 1 {
		return fmt.Sprintf("%d %ss", n, singularNoun(e.Target.Text()))
	}
	if e.Amount != 1 {
		return fmt.Sprintf("%s %d times", e.Target.Text(), e.Amount)
	}
	return e.Target.Text()
}

// foldable folds a single-Æmber exalt, whether it exalts one creature or spreads
// across several; only the repeated "X 2 times" form keeps its own shape.
func (e Exalt) foldable() bool { return e.Amount == 1 }

// Resolve chooses a creature (through the Target) and places Amount Æmber on it.
func (e Exalt) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate exalts each selected creature and reports whether any was exalted,
// so an exalt can be the first half of a Then — Saurus Rex only searches its deck
// when the optional exalt happens. Times repeats the exalt, choosing a fresh
// Target each time, and Distinct spreads those repeats across different creatures.
func (e Exalt) resolveGate(ctx *EffectContext) bool {
	reps := 1
	if e.Times != nil {
		reps = e.Times.Value(ctx)
	}
	var exalted []LocalID
	for i := 0; i < reps; i++ {
		ids := e.Target.selectWith(ctx, false, e.eligible(exalted))
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			exalted = append(exalted, id)
			ctx.Resolver.AddAmberOn(id, e.Amount)
			ctx.Resolver.Record(AemberExalted{
				Player:   ctx.Controller,
				Creature: id,
				Amount:   e.Amount,
			})
		}
	}
	return len(exalted) > 0
}

// eligible bars a creature an earlier repetition already exalted, but only when
// Distinct asks for it. A nil filter leaves the target's candidate set untouched.
func (e Exalt) eligible(exalted []LocalID) func(LocalID) bool {
	if !e.Distinct {
		return nil
	}
	return func(id LocalID) bool { return !slices.Contains(exalted, id) }
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
	ctx.Resolver.Record(AemberExalted{
		Player:   ctx.Controller,
		Creature: id,
		Amount:   e.Amount,
	})
	return true
}
