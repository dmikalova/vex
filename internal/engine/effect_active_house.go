package engine

// MakeItsHouseActive makes the house of the card in context (ctx.It) the active
// player's active house for the rest of the turn, so they may play and use that
// house's cards from now on. Book of leQ pairs it with a preceding RevealTopOfDeck
// so "it" is the revealed top card. It does nothing when no card is in context.
type MakeItsHouseActive struct{}

// Text renders the effect.
func (MakeItsHouseActive) Text() string { return "its house becomes your active house" }

// Resolve sets the active house to the context card's house.
func (MakeItsHouseActive) Resolve(ctx *EffectContext) {
	if !ctx.HasIt {
		return
	}
	ctx.Resolver.SetActiveHouse(ctx.Resolver.House(ctx.It))
}

// SetActiveHouse makes h the active player's active house for the current turn,
// changing which house they may play and use without going through a house choice
// (Book of leQ). It does not re-fire the "after you choose a house" triggers,
// because no house is being chosen — the active house is simply reassigned.
func (g *Game) SetActiveHouse(h House) {
	g.State.ActiveHouse = h
	g.record(HouseChosen{Player: g.State.ActivePlayer, House: h})
}
