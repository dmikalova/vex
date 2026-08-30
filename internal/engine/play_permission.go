package engine

import "fmt"

// PlayPermission is a continuous grant a card makes while in play: its controller
// may play up to Count cards of House on a turn where House is not their active
// house — the off-house play that Witch of the Wilds allows. Splitting the house
// it frees from the count (and, later, a card-type filter or a this-turn window)
// keeps the several off-house-play cards expressible from one shape instead of a
// bespoke field per variant. The zero value (HouseNone) grants nothing.
type PlayPermission struct {
	House House
	Count int
}

// granted reports whether the permission frees any play.
func (p PlayPermission) granted() bool { return p.House != HouseNone }

// count is how many off-house plays the permission allows each turn.
func (p PlayPermission) count() int { return p.Count }

// validate rejects a granted permission that did not state a positive count.
func (p PlayPermission) validate() error {
	if p.granted() && p.Count < 1 {
		return fmt.Errorf("PlayPermission: Count must be positive")
	}
	return nil
}

// WithPlayPermission makes the card, while in play, let its controller play up to
// Count cards of the permission's house on turns where that house is not active.
func WithPlayPermission(p PlayPermission) CardOption {
	return func(c *CardDefinition) { c.PlayPermission = p }
}
