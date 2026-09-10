package engine

import "fmt"

// GainAssault gives each creature its Target selects Assault equal to a Count for
// the remainder of the turn — Creed of Nature grants a chosen creature assault
// equal to its power. The ready phase clears the bonus. It is the lasting-Assault
// companion to GainStats (power and armor). Because it renders a duration body it
// composes under ForDuration with another per-creature grant, so "gains skirmish
// and assault equal to its power" reads under one shared clause.
type GainAssault struct {
	Target Target
	Amount Count
}

// validate requires an explicit target and a magnitude to grant.
func (e GainAssault) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("GainAssault")
	}
	if e.Amount == nil {
		return fmt.Errorf("GainAssault: Amount is required")
	}
	return nil
}

// durationSubject names the creature that gains the Assault, folding the "gains"
// verb into the subject so ForDuration can share it with a sibling keyword grant.
func (e GainAssault) durationSubject() string { return e.Target.Text() + " gains" }

// durationPredicate renders the Assault gained, e.g. "assault equal to its power".
func (e GainAssault) durationPredicate() string {
	return "assault equal to " + e.Amount.CountText()
}

// Text renders the effect, e.g. "for the remainder of the turn, it gains assault
// equal to its power".
func (e GainAssault) Text() string {
	return "for the remainder of the turn, " + e.durationSubject() + " " + e.durationPredicate()
}

// Resolve grants each selected creature the Assault for the remainder of the turn.
// A zero or negative amount grants nothing.
func (e GainAssault) Resolve(ctx *EffectContext) {
	amount := e.Amount.Value(ctx)
	if amount <= 0 {
		return
	}
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.GainAssault(id, amount)
	}
}
