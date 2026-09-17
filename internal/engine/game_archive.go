package engine

// This file holds the archives: taking archived cards into hand at the
// start of a turn, archiving from hand or the top of the deck, and discarding a
// whole archive.

// offerArchives asks a player — after they have chosen their house — whether to
// take their archived cards into their hand, moving the archives to hand if they
// accept. A player with no archived cards is not prompted.
func (g *Game) offerArchives(player int) {
	arc := &g.State.Archives[player]
	if arc.Count == 0 {
		return
	}
	if src, ok := g.constantSelectiveArchivePickup(player); ok {
		g.selectiveArchivePickup(player, src)
		return
	}
	if g.chooseOption(
		player,
		"",
		"Take all the cards from your archives and put them in your hand?",
		[]string{"Yes", "No"},
	) != 0 {
		return
	}
	n := arc.Count
	for _, id := range arc.slice() {
		// Your archives may hold an enemy card, but your hand may not: an abducted card
		// goes to the hand of whoever owns it.
		g.State.Hand[g.owner(id)].add(id)
	}
	*arc = wideList{}
	g.record(ArchivesTakenIntoHand{Player: player, Count: int(n)})
}

// selectiveArchivePickup lets a player take any number of cards from their
// archives into their hand one at a time, declining to stop — the pickup The
// Archivist grants its controller in place of the all-or-nothing offer. Taking
// none is allowed and records nothing. The prompt is attributed to src, the card
// granting the rule.
func (g *Game) selectiveArchivePickup(player int, src LocalID) {
	ctx := &EffectContext{Resolver: g, Controller: player, Source: src}
	chosen := pickCards(
		ctx,
		"Choose a card to take from your archives into your hand",
		0,
		true,
		func() []LocalID { return g.State.Archives[player].slice() },
	)
	for _, id := range chosen {
		g.State.Archives[player].remove(id)
		// Your archives may hold an enemy card, but your hand may not: an abducted card
		// goes to the hand of whoever owns it.
		g.State.Hand[g.owner(id)].add(id)
	}
	if len(chosen) > 0 {
		g.record(ArchivesTakenIntoHand{Player: player, Count: len(chosen)})
	}
}

// archiveFromHand moves a card from a player's hand to their archives.
func (g *Game) archiveFromHand(player int, id LocalID) {
	if g.State.Hand[player].remove(id) {
		g.State.Archives[player].add(id)
		g.record(CardMoved{Player: player, Card: id, From: Hand, To: Archives})
	}
}

// archiveEnemyFromHand moves a card from its owner's hand into a different player's
// archives — an abduction from hand (Hidden Stash). Your archives may hold an enemy
// card, and the ownership rule returns it home the moment it leaves them.
func (g *Game) archiveEnemyFromHand(player int, id LocalID) {
	o := g.owner(id)
	if g.State.Hand[o].remove(id) {
		g.State.Archives[player].add(id)
		g.record(CardAbducted{Player: player, Card: id, Owner: o})
	}
}

func (g *Game) archiveFromDiscard(player int, id LocalID) {
	if g.State.Discard[player].remove(id) {
		g.State.Archives[player].add(id)
		g.record(CardMoved{Player: player, Card: id, From: Discard, To: Archives})
	}
}

// archiveFromPurge moves a card from a player's purge pile to their archives —
// the recovery of a card set aside out of the game (Universal Recycle Bin).
func (g *Game) archiveFromPurge(player int, id LocalID) {
	if g.State.Purge[player].remove(id) {
		g.State.Archives[player].add(id)
		g.record(CardArchivedFromPurge{Player: player, Card: id})
	}
}

// archiveFromDeck moves a specific card the controller looked at — one of the top
// few, not blindly the top one — from a player's deck to their archives.
func (g *Game) archiveFromDeck(player int, id LocalID) {
	if g.State.Deck[player].remove(id) {
		g.State.Archives[player].add(id)
		g.record(TopOfDeckArchived{Player: player, Card: id})
	}
}

// discardArchives moves all of a player's archived cards to their discard pile.
// The active player performs the discard, so they choose the order when it is
// their own archives but cannot when it is an opponent's — those enter the
// discard in a random order, since the active player cannot see them.
func (g *Game) discardArchives(owner int) {
	arc := &g.State.Archives[owner]
	if arc.Count == 0 {
		return
	}
	ids := cloneIDs(arc.slice())
	if owner == g.State.ActivePlayer {
		ids = g.orderByChoice(owner, "Choose the order to discard your archives", ids)
	} else {
		g.State.PRNG.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	}
	*arc = wideList{}
	for _, id := range ids {
		// A discard pile only ever holds its own player's cards, so an abducted card
		// discarded out of these archives goes to its owner's pile.
		g.State.Discard[g.owner(id)].add(id)
	}
	g.record(ArchivesDiscarded{Player: owner, Count: len(ids)})
}

// DiscardCardFromArchives moves a specific card from a player's archives to a
// discard pile, doing nothing if the card is not in those archives. A player's
// archives are facedown, so the card is chosen at random by the caller (the
// Random selection behind DiscardCard{Zones: []Zone{Archives}}), not shown to be picked.
func (g *Game) DiscardCardFromArchives(owner int, id LocalID) {
	arc := &g.State.Archives[owner]
	if arc.indexOf(id) < 0 {
		return
	}
	arc.remove(id)
	// A discard pile only ever holds its own player's cards, so an abducted card
	// discarded out of these archives goes to its owner's pile.
	g.State.Discard[g.owner(id)].add(id)
	g.record(CardMoved{Player: owner, Card: id, From: Archives, To: Discard})
}
