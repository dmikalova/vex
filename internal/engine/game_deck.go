package engine

// This file holds the deck: drawing cards (recycling the discard when the draw
// zone runs out) and shuffling.

// drawOne draws a single card to the player's hand. If the deck is empty, the
// discard pile is shuffled to form a new deck first (KeyForge recycles the
// discard when the draw zone runs out). It returns false only when both the deck
// and discard are empty, so nothing can be drawn.
func (g *Game) drawOne(player int) bool {
	deck := &g.State.Deck[player]
	if deck.Count == 0 {
		if g.State.Discard[player].Count == 0 {
			return false
		}
		g.shuffleDiscardIntoDeck(player)
	}
	g.State.Hand[player].add(deck.removeAt(0))
	return true
}

// shuffleDiscardIntoDeck moves a player's whole discard pile into their deck and
// shuffles it — how KeyForge recycles the discard, both when the draw zone runs
// out and when a card (Help from Future Self) calls for it directly.
func (g *Game) shuffleDiscardIntoDeck(player int) {
	deck, discard := &g.State.Deck[player], &g.State.Discard[player]
	for discard.Count > 0 {
		deck.add(discard.removeAt(0))
	}
	g.Shuffle(player)
}

func (g *Game) shuffleZonesIntoDeck(player int, zones []Zone) {
	deck := &g.State.Deck[player]
	for _, z := range zones {
		switch z {
		case Hand:
			hand := &g.State.Hand[player]
			for _, id := range hand.slice() {
				deck.add(id)
			}
			*hand = deckList{}
		case Archives:
			arc := &g.State.Archives[player]
			for _, id := range arc.slice() {
				deck.add(id)
			}
			*arc = wideList{}
		default: // Discard
			discard := &g.State.Discard[player]
			for _, id := range discard.slice() {
				deck.add(id)
			}
			*discard = deckList{}
		}
	}
	g.Shuffle(player)
}

// draw draws count cards into the player's hand, stopping early only when the
// deck and discard are both exhausted.
func (g *Game) draw(player, count int) {
	for i := 0; i < count; i++ {
		if !g.drawOne(player) {
			break
		}
	}
}

// drawTo draws until the hand holds n cards (or nothing is left to draw).
func (g *Game) drawTo(player, n int) {
	for int(g.State.Hand[player].Count) < n {
		if !g.drawOne(player) {
			break
		}
	}
}

// canDraw reports whether a player still has a card to draw — one in their deck, or
// one in their discard pile that would be reshuffled in.
func (g *Game) canDraw(player int) bool {
	return g.State.Deck[player].Count > 0 || g.State.Discard[player].Count > 0
}

// DiscardTopOfDeck moves the top card of a player's deck to their discard pile.
// Unlike drawing, discarding from the deck does not recycle the discard pile: if
// there is no top card to discard, the effect simply does as much as it can.
func (g *Game) DiscardTopOfDeck(player int) (LocalID, bool) {
	deck := &g.State.Deck[player]
	if deck.Count == 0 {
		return 0, false
	}
	id := deck.removeAt(0)
	g.State.Discard[player].add(id)
	g.record(TopOfDeckDiscarded{Player: player, Card: id})
	return id, true
}

// SwapDeckAndDiscard exchanges a player's deck with their discard pile, then
// shuffles the new deck. Both zones are the same flat list type, so the swap is
// a plain exchange of values (Reverse Time). That exchange is also the right
// ordering: a deck is stored top-first and a discard pile bottom-first, so
// handing one list to the other zone flips the stack over the way turning a
// face-down deck face-up does — the deck's top card becomes the discard's bottom
// card, and the discard's top card becomes the deck's bottom card.
func (g *Game) SwapDeckAndDiscard(player int) {
	deck, discard := &g.State.Deck[player], &g.State.Discard[player]
	*deck, *discard = *discard, *deck
	g.Shuffle(player)
	g.record(DeckAndDiscardSwapped{Player: player})
}

// Shuffle randomizes a player's deck using the game's seeded RNG.
func (g *Game) Shuffle(player int) {
	d := &g.State.Deck[player]
	for i := int(d.Count) - 1; i > 0; i-- {
		j := g.rng.Intn(i + 1)
		d.IDs[i], d.IDs[j] = d.IDs[j], d.IDs[i]
	}
}
