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
	return "for the remainder of the turn, " + e.durationSubject() + " " + e.durationPredicate()
}

// durationSubject and durationPredicate split the body so ForDuration can state
// the shared clause and subject once: "it" / "cannot be dealt damage".
func (e PreventDamage) durationSubject() string   { return e.Target.Text() }
func (e PreventDamage) durationPredicate() string { return "cannot be dealt damage" }

// Resolve marks each selected creature damage-immune for the duration.
func (e PreventDamage) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.PreventDamage(id)
	}
}
