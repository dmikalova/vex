package engine

import "fmt"

// A +1 power counter is a permanent token placed on a creature that raises its
// power by one for as long as it stays in play; a -1 power counter lowers it. A
// creature can hold any number of counters, and they are shed when it leaves play.
//
//rulebook:effect Power Counter

// Power counters are permanent +1/-1 power tokens placed on a creature: they
// raise (or lower) its power for as long as it stays in play, and are shed when
// it leaves. AddPowerCounter places counters on each creature its Target selects.
type AddPowerCounter struct {
	Target Target
	Amount int
}

// validate requires an explicit target.
func (e AddPowerCounter) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("AddPowerCounter")
	}
	return nil
}

// Text renders the effect, e.g. "give Eater of the Dead a +1 power counter".
func (e AddPowerCounter) Text() string {
	return fmt.Sprintf("give %s a %+d power counter", e.Target.Text(), e.Amount)
}

// Resolve places the counters on each selected creature.
func (e AddPowerCounter) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddPowerCounter(id, e.Amount)
	}
}
