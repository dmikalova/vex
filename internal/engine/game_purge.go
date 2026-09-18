package engine

// This file holds the purge pile: cards set aside out of the game. Purging pulls
// a card from a discard pile into its owner's purge pile, from which nothing in
// the base game returns — so purge is the game's way of permanently answering
// recursion out of the discard.

// purgeFromDiscard moves a card from a player's discard pile to its owner's purge
// pile. Callers pass a card already in that discard pile.
func (g *Game) purgeFromDiscard(holder int, id LocalID) {
	g.purgeFrom(holder, id, Discard)
}

// purgeFromHand moves a card from a player's hand to its owner's purge pile. The
// card whose ability purged it is credited through the record's frame. Callers
// pass a card already in that hand.
func (g *Game) purgeFromHand(holder int, id LocalID) {
	g.moveCard(id,
		zoneRef{Player: holder, Zone: Hand},
		zoneRef{Player: g.owner(id), Zone: Purged},
		CardPurgedFromHand{Card: id, Owner: holder})
}

// purgeFromArchives moves a card from a player's archives to its owner's purge
// pile. Callers pass a card already in those archives.
func (g *Game) purgeFromArchives(holder int, id LocalID) {
	g.purgeFrom(holder, id, Archives)
}

// purgeFromDeck moves a card from a player's deck to its owner's purge pile.
// Callers pass a card already in that deck (Borr Nit purges one of the cards it
// revealed off the top).
func (g *Game) purgeFromDeck(holder int, id LocalID) {
	g.purgeFrom(holder, id, Deck)
}

// purgeFrom purges from a zone whose move narrates as a plain relocation. Purging
// from hand is the exception and records its own entry: a hand is hidden, so the
// log has to say the card came from one.
//
// holder is who held the source zone, which is not always the owner: archives can
// hold an abducted enemy card (Hidden Stash), and a purged card always goes to its
// owner's pile. Pinned by TestPurgeFromArchivesGoesToOwnersPile.
func (g *Game) purgeFrom(holder int, id LocalID, from Zone) {
	g.moveCard(id,
		zoneRef{Player: holder, Zone: from},
		zoneRef{Player: g.owner(id), Zone: Purged},
		CardMoved{Player: g.State.ActivePlayer, Card: id, From: from, To: Purged})
}
