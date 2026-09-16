package engine

import "fmt"

// GainTrait gives each creature its Target selects a trait until the start of the
// controller's next turn — the Mutation cycle grants a chosen creature the Mutant
// trait. Only the controller's own start-of-turn phase clears it, before any
// start-of-turn ability resolves, so the trait survives the opponent's turn, where
// an enemy card may read it. Because it renders a duration
// body it folds under GainUntilNextTurn beside a keyword or Assault grant, so "gains
// skirmish and the Mutant trait" reads under one shared clause.
type GainTrait struct {
	Target Target
	Trait  Trait
}

// validate requires an explicit target and a real trait.
func (e GainTrait) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("GainTrait")
	}
	if e.Trait == traitUnset {
		return fmt.Errorf("GainTrait: unset trait")
	}
	return nil
}

// durationSubject names the creature that gains the trait, folding the "gains"
// verb into the subject so GainUntilNextTurn can share it with a sibling grant.
func (e GainTrait) durationSubject() string { return e.Target.Text() + " gains" }

// durationPredicate renders the trait gained, e.g. "the Mutant trait".
func (e GainTrait) durationPredicate() string { return "the " + e.Trait.String() + " trait" }

// Text renders the effect, e.g. "it gains the Mutant trait until the start of your
// next turn".
func (e GainTrait) Text() string {
	return e.durationSubject() + " " + e.durationPredicate() +
		" until the start of your next turn"
}

// Resolve grants each selected creature the trait until the controller's next turn.
func (e GainTrait) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.GrantTraitUntilNextTurn(id, e.Trait)
	}
}
