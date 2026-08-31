package engine

import (
	"fmt"
	"strings"
)

// ShuffleIntoDeck shuffles the controller's named Zones into their deck — the
// discard pile (Help from Future Self), the hand and discard pile (Screaming
// Cave), or the archives and discard pile. It moves every named zone's cards into
// the deck, then shuffles once.
//
//rulebook:effect Shuffle Into Deck
type ShuffleIntoDeck struct {
	Zones []Zone
}

// validate requires at least one shuffleable zone.
func (e ShuffleIntoDeck) validate() error {
	if len(e.Zones) == 0 {
		return fmt.Errorf("ShuffleIntoDeck: at least one zone must be set")
	}
	for _, z := range e.Zones {
		if !shuffleableZone(z) {
			return fmt.Errorf("ShuffleIntoDeck: zone %d cannot be shuffled into the deck", z)
		}
	}
	return nil
}

// Text renders the effect, e.g. "shuffle your hand and discard pile into your deck".
func (e ShuffleIntoDeck) Text() string {
	nouns := make([]string, len(e.Zones))
	for i, z := range e.Zones {
		nouns[i] = shuffleZoneNoun(z)
	}
	return "shuffle your " + strings.Join(nouns, " and ") + " into your deck"
}

// Resolve moves the controller's named zones into their deck and shuffles.
func (e ShuffleIntoDeck) Resolve(ctx *EffectContext) {
	ctx.Resolver.ShuffleZonesIntoDeck(ctx.Controller, e.Zones)
}

// shuffleableZone reports whether a zone can be shuffled into the deck.
func shuffleableZone(z Zone) bool {
	return z == Discard || z == Hand || z == Archives
}

// shuffleZoneNoun names a zone for the shuffle phrase.
func shuffleZoneNoun(z Zone) string {
	switch z {
	case Hand:
		return "hand"
	case Archives:
		return "archives"
	default: // Discard
		return "discard pile"
	}
}
