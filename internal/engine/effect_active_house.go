package engine

import "fmt"

// ChangeActiveHouse changes the active player's active house for the rest of the
// turn to the house To names, so they may play and use that house's cards from now
// on. Book of leQ sets it to the revealed card's house (TheContextualHouse), paired
// with a preceding RevealTopOfDeck so "it" is the revealed top card. It does nothing
// when To names the card in context and none is present.
type ChangeActiveHouse struct {
	// To names the house to switch to — TheContextualHouse for the card in context
	// (Book of leQ), TheChosenHouse for a house an enclosing ChooseHouseThen picked.
	To HouseChoice
}

// validate requires a house source.
func (e ChangeActiveHouse) validate() error {
	if e.To == houseChoiceUnset {
		return fmt.Errorf("ChangeActiveHouse: To must be set")
	}
	return nil
}

// Text renders the effect, e.g. "its house becomes your active house".
func (e ChangeActiveHouse) Text() string {
	switch e.To {
	case TheChosenHouse:
		return "the chosen house becomes your active house"
	default:
		return "its house becomes your active house"
	}
}

// Resolve sets the active house to the house To names, leaving it unchanged when
// that house is not available (no card in context).
func (e ChangeActiveHouse) Resolve(ctx *EffectContext) {
	h := e.To.resolveHouse(ctx)
	if h == HouseNone {
		return
	}
	ctx.Resolver.SetActiveHouse(h)
}

// SetActiveHouse makes h the active player's active house for the current turn,
// changing which house they may play and use without going through a house choice
// (Book of leQ). It does not re-fire the "after you choose a house" triggers,
// because no house is being chosen — the active house is simply reassigned.
func (g *Game) SetActiveHouse(h House) {
	g.State.ActiveHouse = h
	g.record(HouseChosen{Player: g.State.ActivePlayer, House: h})
}
