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
	// zone is the zone half of the destination.
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
	// The zones a Destination can name.
	destHand
	destTopOfDeck
	destBottomOfDeck
	destDeckShuffled
	destArchives
	// The two terminal destinations a card leaves the game or is set aside in.
	// They are the endpoints of the purge and discard verbs, so no authoring To
	// value names them: a card writes Purge/Discard, and those verbs move through
	// the unexported toPurged/toDiscard destinations below (ADR 0031).
	destDiscard
	destPurged
)

// The destinations a card can name.
var (
	ToHand         = Destination{zone: destHand}
	ToTopOfDeck    = Destination{zone: destTopOfDeck}
	ToBottomOfDeck = Destination{zone: destBottomOfDeck}
	ToDeckShuffled = Destination{zone: destDeckShuffled}
	ToArchives     = Destination{zone: destArchives}
)

// The terminal destinations, unexported because a card names the verb (Purge,
// Discard), not the destination: the purge and discard effects move through these
// so their dispatch shares the one source-by-destination matrix (ADR 0031).
var (
	toDiscard = Destination{zone: destDiscard}
	toPurged  = Destination{zone: destPurged}
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

// move carries cards out of play to the destination: the from-play arm of the
// matrix, where the resolving player owns whichever pile they land in. It is
// plural by default, because moving several cards is one moment — the board
// settles once after the whole set rather than between its members, so a card
// leaving play cannot destroy one still waiting its turn. Nesting is harmless:
// an inner batch's settle is a no-op while an outer one is open.
func (d Destination) move(ctx *EffectContext, ids ...LocalID) {
	ctx.Resolver.Simultaneously(ctx.Controller, func() {
		for _, id := range ids {
			d.moveFrom(ctx, InPlay, ctx.Controller, id)
		}
	})
}

// moveSource carries the card whose ability is resolving to the destination,
// wherever that card is. A source still in play — a creature or artifact — moves
// out of play. A resolving Tactic is in no zone at all, having already left the
// one it was played from, so it cannot be moved; it is instead redirected to end
// its play here rather than in the discard pile. That redirect names no source
// zone, which is why it works whatever the Tactic was played from
// (TestResolvingCardRedirectIsPerCard).
func (d Destination) moveSource(ctx *EffectContext) {
	if resolverInPlay(ctx, ctx.Source) {
		d.move(ctx, ctx.Source)
		return
	}
	ctx.Resolver.RedirectResolvingCard(ctx.Source, d)
}

// moveFrom carries one card from the zone `from` to this destination, dispatching
// to the resolver removal for that source-and-destination pair — the one
// source-by-destination move matrix behind the movement verbs (ADR 0031). owner
// names whose copy of a pile zone the card sits in; the from-play resolvers ignore
// it because a card in play has no owning pile. Only the pairs a verb actually
// produces are listed: purge from a hand, a discard pile, or play; discard from a
// hand or archives; archive from a hand, a discard pile, a deck, or play; into a
// hand from a deck, a discard pile, or play; shuffled into a deck from a hand, a
// discard pile, or play; and on top of a deck from play or a discard pile. archive
// also has a purge-pile source (Universal Recycle Bin recovering a card set aside
// for good).
func (d Destination) moveFrom(ctx *EffectContext, from Zone, owner int, id LocalID) {
	switch d.zone {
	case destPurged:
		switch from {
		case Hand:
			ctx.Resolver.PurgeFromHand(owner, id)
		case Discard:
			ctx.Resolver.PurgeFromDiscard(owner, id)
		case Archives:
			ctx.Resolver.PurgeFromArchives(owner, id)
		case Deck:
			ctx.Resolver.PurgeFromDeck(owner, id)
		default: // inPlay
			ctx.Resolver.PurgeFromPlay(id)
		}
	case destDiscard:
		if from == Archives {
			ctx.Resolver.DiscardCardFromArchives(owner, id)
			return
		}
		ctx.Resolver.DiscardCardFromHand(owner, id)
	case destTopOfDeck:
		switch from {
		case Discard:
			ctx.Resolver.MoveFromDiscardToTopOfDeck(id)
		case Deck:
			ctx.Resolver.MoveFromDeckToTopOfDeck(id)
		default:
			ctx.Resolver.PutOnTopOfDeck(id)
		}
	case destDeckShuffled:
		switch from {
		case Hand:
			ctx.Resolver.ShuffleFromHandIntoDeck(id)
		case Discard:
			ctx.Resolver.ShuffleFromDiscardIntoDeck(id)
		default: // inPlay
			ctx.Resolver.PutIntoDeckShuffled(id)
		}
	case destArchives:
		switch from {
		case Hand:
			ctx.Resolver.ArchiveFromHand(id)
		case Discard:
			ctx.Resolver.ArchiveFromDiscard(owner, id)
		case Deck:
			ctx.Resolver.ArchiveFromDeck(id)
		case Purged:
			ctx.Resolver.ArchiveFromPurge(owner, id)
		default: // inPlay
			if d.yours {
				ctx.Resolver.PutIntoYourArchives(id, owner)
				return
			}
			ctx.Resolver.PutIntoArchives(id)
		}
	default: // destHand
		switch from {
		case Deck:
			ctx.Resolver.MoveFromDeckToHand(id)
		case Discard:
			ctx.Resolver.PutFromDiscardIntoHand(id)
		default: // inPlay
			ctx.Resolver.PutIntoHand(id)
		}
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
