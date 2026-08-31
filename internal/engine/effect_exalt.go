package engine

import "fmt"

// To exalt a creature is to place 1 Aember from the common supply onto a chosen
// friendly or enemy creature. The Aember sits on the creature, belonging to no
// pool, until it leaves play, then goes to the owner's opponent's pool. Exalting
// N times places N Aember.
//
//rulebook:effect Exalt
type Exalt struct {
	Target Target
	Times  int
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
	if e.Times == 1 {
		return "exalt " + e.Target.Text()
	}
	return fmt.Sprintf("exalt %s %d times", e.Target.Text(), e.Times)
}

// Resolve chooses a creature (through the Target) and places Times Æmber on it.
func (e Exalt) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddAmberOn(id, e.Times)
		ctx.Resolver.Logf("%s is exalted (%d Æmber placed)", ctx.Resolver.Name(id), e.Times)
	}
}
