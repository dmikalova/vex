package engine

// This file holds the purge pile: cards set aside out of the game. Purging pulls
// a card from a discard pile into its owner's purge pile, from which nothing in
// the base game returns — so purge is the game's way of permanently answering
// recursion out of the discard.

// purgeFromDiscard moves a card from a player's discard pile to their purge pile.
// Callers pass a card already in that discard pile.
func (g *Game) purgeFromDiscard(owner int, id LocalID) {
	g.purgeFrom(owner, id, Discard)
}

// purgeFromHand moves a card from a player's hand to their purge pile. The card
// whose ability purged it is credited through the record's frame. Callers pass a
// card already in that hand.
func (g *Game) purgeFromHand(owner int, id LocalID) {
	g.moveCard(id,
		zoneRef{Player: owner, Zone: Hand},
		zoneRef{Player: owner, Zone: purged},
		CardPurgedFromHand{Card: id, Owner: owner})
}

// purgeFromArchives moves a card from a player's archives to their purge pile.
// Callers pass a card already in those archives.
func (g *Game) purgeFromArchives(owner int, id LocalID) {
	g.purgeFrom(owner, id, Archives)
}

// purgeFromDeck moves a card from a player's deck to their purge pile. Callers
// pass a card already in that deck (Borr Nit purges one of the cards it revealed
// off the top).
func (g *Game) purgeFromDeck(owner int, id LocalID) {
	g.purgeFrom(owner, id, Deck)
}

// purgeFrom purges from a zone whose move narrates as a plain relocation. Purging
// from hand is the exception and records its own entry: a hand is hidden, so the
// log has to say the card came from one.
func (g *Game) purgeFrom(owner int, id LocalID, from Zone) {
	g.moveCard(id,
		zoneRef{Player: owner, Zone: from},
		zoneRef{Player: owner, Zone: purged},
		CardMoved{Player: g.State.ActivePlayer, Card: id, From: from, To: purged})
}
