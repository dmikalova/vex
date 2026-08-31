package engine

import "fmt"

// BelongToHouse makes each creature its Target selects belong to House for the
// given Duration, overriding the house it counts as for active-house checks (Brain
// Stem Antenna's host counts as Mars for the rest of the turn). The change is
// per-match state, dropped when the creature leaves play; EndOfTurn also drops it at
// end of turn, while UntilThisLeavesPlay keeps it until the creature leaves play.
//
//rulebook:effect Belong to House
type BelongToHouse struct {
	Target   Target
	House    House
	Duration Duration
}

// validate requires a target, a house, and a duration this effect supports.
func (e BelongToHouse) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("BelongToHouse")
	}
	if e.House == HouseNone {
		return fmt.Errorf("BelongToHouse: house must be set")
	}
	if e.Duration != EndOfTurn && e.Duration != UntilThisLeavesPlay {
		return fmt.Errorf("BelongToHouse: duration must be EndOfTurn or UntilThisLeavesPlay")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn this creature
// belongs to house Mars".
func (e BelongToHouse) Text() string {
	if e.Duration == UntilThisLeavesPlay {
		return e.Target.Text() + " belongs to house " + e.House.String() + " until it leaves play"
	}
	return "for the remainder of the turn " + e.Target.Text() + " belongs to house " + e.House.String()
}

// Resolve makes each selected creature belong to House for the Duration.
func (e BelongToHouse) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if e.Duration == UntilThisLeavesPlay {
			ctx.Resolver.SetLastingHouse(id, e.House)
		} else {
			ctx.Resolver.BelongToHouseForRemainderOfTurn(id, e.House)
		}
	}
}
