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

// archiveFromHand moves a card from a player's hand to their archives.
func (g *Game) archiveFromHand(player int, id LocalID) {
	if g.State.Hand[player].remove(id) {
		g.State.Archives[player].add(id)
		g.record(CardArchivedFromHand{Player: player, Card: id})
	}
}

func (g *Game) archiveFromDiscard(player int, id LocalID) {
	if g.State.Discard[player].remove(id) {
		g.State.Archives[player].add(id)
		g.record(CardArchivedFromDiscard{Player: player, Card: id})
	}
}

// archiveTopOfDeck moves the top card of a player's deck to their archives,
// reporting whether a card was available to archive.
func (g *Game) archiveTopOfDeck(player int) bool {
	deck := &g.State.Deck[player]
	if deck.Count == 0 {
		return false
	}
	id := deck.removeAt(0)
	g.State.Archives[player].add(id)
	g.record(TopOfDeckArchived{Player: player, Card: id})
	return true
}

// archiveTopOfDiscard moves the top card of a player's discard pile — the most
// recently discarded, at the end of the list — to their archives, reporting
// whether a card was available to archive.
func (g *Game) archiveTopOfDiscard(player int) bool {
	discard := &g.State.Discard[player]
	if discard.Count == 0 {
		return false
	}
	id := discard.removeAt(int(discard.Count) - 1)
	g.State.Archives[player].add(id)
	g.record(CardArchivedFromDiscard{Player: player, Card: id})
	return true
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
		g.rng.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	}
	*arc = wideList{}
	for _, id := range ids {
		// A discard pile only ever holds its own player's cards, so an abducted card
		// discarded out of these archives goes to its owner's pile.
		g.State.Discard[g.owner(id)].add(id)
	}
	g.record(ArchivesDiscarded{Player: owner, Count: len(ids)})
}
