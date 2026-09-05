package engine

import (
	"fmt"
	"strings"
)

// ShuffleIntoDeck shuffles the controller's named Zones into their deck — the
// discard pile (Help from Future Self), the hand and discard pile (Screaming
// Cave), or the archives and discard pile. It moves every named zone's cards into
// the deck, then shuffles once.
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
		nouns[i] = z.noun()
	}
	return "shuffle your " + strings.Join(nouns, " and ") + " into your deck"
}

// Resolve moves the controller's named zones into their deck and shuffles.
func (e ShuffleIntoDeck) Resolve(ctx *EffectContext) {
	ctx.Resolver.ShuffleZonesIntoDeck(ctx.Controller, e.Zones)
}

// SwapDeckAndDiscard exchanges the controller's deck with their discard pile and
// shuffles the new deck — Reverse Time turns a spent deck back into a fresh one.
// It differs from ShuffleIntoDeck{Discard} in that the old deck goes away into
// the discard pile rather than surviving underneath it.
type SwapDeckAndDiscard struct{}

// Text renders the effect.
func (e SwapDeckAndDiscard) Text() string {
	return "swap your deck and your discard pile, then shuffle your deck"
}

// Resolve swaps the two zones and shuffles.
func (e SwapDeckAndDiscard) Resolve(ctx *EffectContext) {
	ctx.Resolver.SwapDeckAndDiscard(ctx.Controller)
}

// shuffleableZone reports whether a zone can be shuffled into the deck.
func shuffleableZone(z Zone) bool {
	return z == Discard || z == Hand || z == Archives
}
