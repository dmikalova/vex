package engine

import "fmt"

// Destination names where an effect puts a card it moves, and renders the
// KeyForge phrase for it — "into its owner's hand", "on top of its owner's deck",
// "into your archives". It is a zone plus whose copy of that zone: by the
// ownership rule only your battleline, artifact line, and archives may hold an
// enemy card, so Yours is meaningful on the archives alone. The deck has three
// separate destinations (its top, its bottom, and shuffled in), so "the deck"
// alone is never a destination; a card always names which.
type Destination struct {
	zone destinationZone
	// yours redirects the move into the resolving player's own copy of the zone.
	yours bool
}

// destinationZone is the zone half of a Destination.
type destinationZone uint8

const (
	// destUnset is the invalid zero value: an effect must name where a card goes
	// rather than leave the destination unset.
	destUnset destinationZone = iota
	destHand
	destTopOfDeck
	destBottomOfDeck
	destDeckShuffled
	destArchives
)

// The destinations a card can name.
var (
	ToHand         = Destination{zone: destHand}
	ToTopOfDeck    = Destination{zone: destTopOfDeck}
	ToBottomOfDeck = Destination{zone: destBottomOfDeck}
	ToDeckShuffled = Destination{zone: destDeckShuffled}
	ToArchives     = Destination{zone: destArchives}
)

// Yours sends the card to the resolving player's own copy of the zone rather than
// its owner's — an abduction into your archives (Uxlyx the Zookeeper).
func (d Destination) Yours() Destination {
	d.yours = true
	return d
}

// destinationClause is the phrase each zone renders, singular then plural. The
// verb is part of it: shuffling a card into a deck is not a "put".
var destinationClause = map[destinationZone][2]string{
	destHand: {
		"put %s into its owner's hand",
		"put %s into their owners' hands",
	},
	destTopOfDeck: {
		"put %s on top of its owner's deck",
		"put %s on top of their owners' decks",
	},
	destBottomOfDeck: {
		"put %s on the bottom of its owner's deck",
		"put %s on the bottom of their owners' decks",
	},
	destDeckShuffled: {
		"shuffle %s into its owner's deck",
		"shuffle %s into their owners' decks",
	},
	destArchives: {
		"put %s into its owner's archives",
		"put %s into their owners' archives",
	},
}

// clause renders the whole move for the cards named by subject, e.g.
// "put an enemy creature into your archives".
func (d Destination) clause(subject string, plural bool) string {
	if d.yours {
		return fmt.Sprintf("put %s into your archives", subject)
	}
	form := destinationClause[d.zone]
	if plural {
		return fmt.Sprintf(form[1], subject)
	}
	return fmt.Sprintf(form[0], subject)
}

// move carries one card out of play to the destination.
func (d Destination) move(ctx *EffectContext, id LocalID) {
	switch d.zone {
	case destTopOfDeck:
		ctx.Resolver.PutOnTopOfDeck(id)
	case destDeckShuffled:
		ctx.Resolver.PutIntoDeckShuffled(id)
	case destArchives:
		if d.yours {
			ctx.Resolver.PutIntoYourArchives(id, ctx.Controller)
			return
		}
		ctx.Resolver.PutIntoArchives(id)
	default:
		ctx.Resolver.PutIntoHand(id)
	}
}

// movable reports whether an effect taking a card out of play can send it here.
// The bottom of the deck is not among them: no card puts a card from play there.
func (d Destination) movable() bool {
	switch d.zone {
	case destHand, destTopOfDeck, destDeckShuffled, destArchives:
		return true
	default:
		return false
	}
}
