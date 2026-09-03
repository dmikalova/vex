package engine

import "fmt"

// PreventDamage marks each creature its Target selects immune to damage for the
// given Duration — Shield of Justice protects each friendly creature for the
// remainder of the turn. Only EndOfTurn is supported today (the mark is cleared by
// the ready phase); the field is explicit so a longer window can be added without
// reshaping the callers.
type PreventDamage struct {
	Target   Target
	Duration Duration
}

// validate requires an explicit target and a supported duration.
func (e PreventDamage) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PreventDamage")
	}
	if e.Duration != EndOfTurn {
		return fmt.Errorf("PreventDamage: duration must be EndOfTurn")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, each friendly
// creature cannot be dealt damage".
func (e PreventDamage) Text() string {
	return fmt.Sprintf("for the remainder of the turn, %s cannot be dealt damage", e.Target.Text())
}

// Resolve marks each selected creature damage-immune for the duration.
func (e PreventDamage) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.PreventDamage(id)
	}
}
