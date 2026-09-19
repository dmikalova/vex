package engine

// This file holds the one primitive every resting-zone move goes through. Moving
// a card between two of the piles a player holds — hand, deck, discard, archives,
// purge — is always the same three steps: take it out of one pile, put it into
// another, and narrate the move. The verbs above it (archive, discard, purge, put)
// differ only in which piles they name and which entry they record, so they pass
// those in rather than each repeating the steps.
//
// Play is not a resting zone and is deliberately absent: a card in play carries
// upgrades, damage, and ward, so entering and leaving it runs whole lifecycles
// that this primitive does not and should not know about. It also has no pile to
// take the card out of, so a play source would leave the `from` side a branch
// whose one arm ignores `from` entirely. Leaving play goes through
// leavePlayInto/fileFromPlay instead, and the two meet one level up in
// Destination.moveFrom, which is the real source axis (ADR 0031).

// cardPile is the pair of operations a resting-zone move needs. Archives are a
// wideList while the other piles are deckLists, so the mover reaches all five
// through the methods they share rather than through their concrete types.
type cardPile interface {
	add(id LocalID)
	remove(id LocalID) bool
}

// zoneRef names one player's copy of a resting zone. A move needs the player as
// well as the zone because a pile can receive a card its owner does not hold —
// archives take abducted enemy cards (Hidden Stash).
type zoneRef struct {
	Player int
	Zone   Zone
}

// pileZones lists every zone pile returns a pile for, so a caller that must
// touch all of them (removeFromRestingZones) enumerates them here rather than
// naming the five state fields again and silently missing a sixth.
var pileZones = [...]Zone{Hand, Deck, Discard, Archives, Purged}

// removeFromRestingZones unlists id from every pile its owner holds and returns
// that owner. A card in play sits in none of them, so a caller that may be handed
// one deals with play itself (manualRelocate does; putIntoPlay guards against it).
func (g *Game) removeFromRestingZones(id LocalID) int {
	o := g.owner(id)
	for _, z := range pileZones {
		g.pile(zoneRef{Player: o, Zone: z}).remove(id)
	}
	return o
}

// pile returns the player's pile for a resting zone, or nil for a zone that holds
// no pile of cards.
func (g *Game) pile(ref zoneRef) cardPile {
	switch ref.Zone {
	case Hand:
		return &g.State.Hand[ref.Player]
	case Deck:
		return &g.State.Deck[ref.Player]
	case Discard:
		return &g.State.Discard[ref.Player]
	case Archives:
		return &g.State.Archives[ref.Player]
	case Purged:
		return &g.State.Purge[ref.Player]
	default:
		return nil
	}
}

// moveCard takes a card out of one resting zone, puts it into another, and records
// entry. It reports whether the card was there to move, so a caller that is not
// certain of the card's whereabouts can gate on the move having happened rather
// than probing the zone itself.
func (g *Game) moveCard(id LocalID, from, to zoneRef, entry LogEntry) bool {
	src, dst := g.pile(from), g.pile(to)
	if src == nil || dst == nil || !src.remove(id) {
		return false
	}
	dst.add(id)
	g.record(entry)
	return true
}
