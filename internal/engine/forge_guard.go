package engine

// opponentForgeGuarded resolves every forge guard the non-forging player controls
// against forger's imminent key forge, reporting whether one of them prevents it.
// It is Keyforgery's interrupt: when the opponent would forge a key, that player
// names a house, a random card is revealed from the guard's controller's hand,
// and if that card is not of the named house the guard is destroyed and the forge
// is prevented (leaving forger's Æmber unspent). A guard whose controller has an
// empty hand cannot act, and a correct guess leaves the guard in play and lets the
// forge proceed.
func (g *Game) opponentForgeGuarded(forger int) bool {
	guard := 1 - forger
	for _, id := range g.allInPlay(guard) {
		if !g.cat.def(id).GuardsOpponentForge {
			continue
		}
		revealed, ok := g.randomCardFromHand(guard)
		if !ok {
			continue
		}
		options := houseNames[1:] // every house except HouseNone
		named := House(g.ChooseOption(forger, id, "Name a house", options) + 1)
		g.record(CardsRevealedToAll{Player: guard, Cards: []LocalID{revealed}})
		if g.House(revealed) != named {
			g.record(KeyForgePrevented{Player: forger, By: id})
			g.destroyEach(guard, []LocalID{id})
			return true
		}
	}
	return false
}
