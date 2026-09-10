package engine

// This file holds the purge pile: cards set aside out of the game. Purging pulls
// a card from a discard pile into its owner's purge pile, from which nothing in
// the base game returns — so purge is the game's way of permanently answering
// recursion out of the discard.

// purgeFromDiscard moves a card from a player's discard pile to their purge pile.
// Callers pass a card already in that discard pile.
func (g *Game) purgeFromDiscard(owner int, id LocalID) {
	g.State.Discard[owner].remove(id)
	g.State.Purge[owner].add(id)
	g.record(CardMoved{Player: g.State.ActivePlayer, Card: id, From: Discard, To: purged})
}

// purgeFromHand moves a card from a player's hand to their purge pile. Callers
// pass a card already in that hand.
func (g *Game) purgeFromHand(owner int, id LocalID) {
	g.State.Hand[owner].remove(id)
	g.State.Purge[owner].add(id)
	g.record(CardMoved{Player: g.State.ActivePlayer, Card: id, From: Hand, To: purged})
}

// purgeFromArchives moves a card from a player's archives to their purge pile.
// Callers pass a card already in those archives.
func (g *Game) purgeFromArchives(owner int, id LocalID) {
	g.State.Archives[owner].remove(id)
	g.State.Purge[owner].add(id)
	g.record(CardMoved{Player: g.State.ActivePlayer, Card: id, From: Archives, To: purged})
}

// purgeFromDeck moves a card from a player's deck to their purge pile. Callers
// pass a card already in that deck (Borr Nit purges one of the cards it revealed
// off the top).
func (g *Game) purgeFromDeck(owner int, id LocalID) {
	g.State.Deck[owner].remove(id)
	g.State.Purge[owner].add(id)
	g.record(CardMoved{Player: g.State.ActivePlayer, Card: id, From: Deck, To: purged})
}
