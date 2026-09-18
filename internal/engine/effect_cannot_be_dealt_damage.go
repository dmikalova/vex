package engine

import "fmt"

// CannotBeDealtDamage marks each creature its Target selects unable to be dealt
// damage for the given Duration — Shield of Justice protects each friendly
// creature for the remainder of the turn, Lucky Dice during the opponent's next
// turn. It installs a continuous effect read live at damage time.
type CannotBeDealtDamage struct {
	Target   Target
	Duration Duration
}

// validate requires an explicit target and one of the supported durations.
func (e CannotBeDealtDamage) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("CannotBeDealtDamage")
	}
	switch e.Duration {
	case RemainderOfPlayerTurn, OpponentNextTurn:
		return nil
	default:
		return fmt.Errorf(
			"CannotBeDealtDamage: duration must be RemainderOfPlayerTurn or OpponentNextTurn")
	}
}

// Text renders the effect, e.g. "for the remainder of the turn, each friendly
// creature cannot be dealt damage".
func (e CannotBeDealtDamage) Text() string {
	return durationClause(e.Duration, "") + ", " + e.durationSubject() + " " + e.durationPredicate()
}

// durationSubject and durationPredicate split the body so ForDuration can state
// the shared clause and subject once: "it" / "cannot be dealt damage".
func (e CannotBeDealtDamage) durationSubject() string   { return e.Target.Text() }
func (e CannotBeDealtDamage) durationPredicate() string { return "cannot be dealt damage" }

// Resolve marks each selected creature unable to be dealt damage for the duration.
// When the target is a whole, unfiltered side (Shield of Justice's "each friendly
// creature") the immunity is registered side-wide and read live, so a creature
// played or gained later is protected too, rather than a snapshot of the
// creatures in play the moment it resolved.
func (e CannotBeDealtDamage) Resolve(ctx *EffectContext) {
	if player, ok := e.Target.wholeSide(ctx.Controller); ok {
		ctx.Resolver.SetSideDamageImmune(player, e.Duration)
		return
	}
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.SetDamageImmune(id, e.Duration)
	}
}
