package engine

import "fmt"

// PlayPermission is a continuous grant a card makes while in play: its controller
// may play up to Amount cards of House on a turn where House is not their active
// house — the off-house play that Witch of the Wilds allows. Splitting the house
// it frees from the count (and, later, a card-type filter or a this-turn window)
// keeps the several off-house-play cards expressible from one shape instead of a
// bespoke field per variant. The zero value (HouseNone) grants nothing.
//
// NonActive frees any house that is not the active house rather than one named
// House, and Condition gates the whole grant behind a live predicate read from the
// source's point of view — together they express Captain Val Jericho, whose
// SourceInCenterOfBattleline condition frees one non-active-house play only while
// it is centered.
type PlayPermission struct {
	House     House
	Amount    int
	NonActive bool
	Condition Condition
	// Types, when set, makes the grant a house-agnostic, unlimited waiver for cards
	// of those types: Matter Maker lets its controller play any number of upgrades
	// as if they were of the active house. It is a mode of its own — the counted
	// House and NonActive axes do not apply — so the play path never consumes it.
	Types CardTypes
}

// granted reports whether the permission frees any play.
func (p PlayPermission) granted() bool {
	return p.House != HouseNone || p.NonActive || p.Types != 0
}

// count is how many off-house plays the permission allows each turn.
func (p PlayPermission) count() int { return p.Amount }

// validate rejects a granted permission that did not state a positive count, or
// that carries a misconfigured condition. A Types-scoped waiver is unlimited, so
// it needs no count.
func (p PlayPermission) validate() error {
	if p.Types == 0 && p.granted() && p.Amount < 1 {
		return fmt.Errorf("PlayPermission: Amount must be positive")
	}
	if p.Condition != nil {
		return validateCondition(p.Condition)
	}
	return nil
}

// WithPlayPermission makes the card, while in play, let its controller play up to
// Amount cards of the permission's house on turns where that house is not active.
func WithPlayPermission(p PlayPermission) CardOption {
	return func(c *CardDefinition) { c.PlayPermission = p }
}
