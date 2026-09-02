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

// NameHouse remembers the house a surrounding ChooseHouseThen picked on the source
// card, where it stays for as long as that card is in play. It is the writer half
// of a HouseLock whose house is not printed but named: Restringuntus chooses a
// house on play and bars its opponent from it until it leaves play. Player names
// whose choice the lock will constrain, and must match the card's HouseLock.
//
//rulebook:effect Name a House
type NameHouse struct {
	Player Player
}

// validate requires the player whose house choice the named house will constrain.
func (e NameHouse) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("NameHouse")
	}
	return nil
}

// Text renders the lock the named house creates, e.g. "your opponent cannot choose
// that house as their active house until {self} leaves play". The house itself is
// named by the enclosing ChooseHouseThen.
func (e NameHouse) Text() string {
	who, possessive := "you", "your"
	if e.Player == Opponent {
		who, possessive = "your opponent", "their"
	}
	return who + " cannot choose that house as " + possessive +
		" active house until " + SelfName + " leaves play"
}

// Resolve stores the chosen house on the source card.
func (NameHouse) Resolve(ctx *EffectContext) {
	ctx.Resolver.SetNamedHouse(ctx.Source, ctx.ChosenHouse)
}
