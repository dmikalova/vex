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

// Resolve chooses a creature (through the Target) and places Amount Æmber on it.
func (e Exalt) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddAmberOn(id, e.Amount)
		ctx.Resolver.Record(AemberExalted{Creature: id, Amount: e.Amount})
	}
}
