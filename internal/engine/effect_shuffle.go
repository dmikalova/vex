package engine

import (
	"fmt"
	"strings"
)

// Shuffle shuffles the controller's deck, optionally folding whole zones of their
// cards into it first (ADR 0031). With no source it is the plain "shuffle your
// deck" a deck search ends on (Orb of Wonder, Saurus Rex, Grumpus Tamer) — a
// standalone effect so a search never bundles its own shuffle, enforced by a card
// lint (TestSearchIsFollowedByShuffle). Zones folds whole piles in before the
// shuffle — the discard pile (Help from Future Self), the hand and discard pile
// (Screaming Cave), the archives and discard pile. To fold every friendly card in
// play into the deck instead, use ShuffleIntoDeck with From: []Zone{InPlay}.
type Shuffle struct {
	// Zones names whole piles that fold into the deck before it shuffles; each must
	// be shuffleable (a hand, a discard pile, or archives). Empty is the bare deck
	// shuffle.
	Zones []Zone
}

// validate rejects an unshuffleable zone. The zero value (a bare deck shuffle) is
// valid, so it also serves as the deck-dig "shuffle that deck" routing terminal
// (effect_deck.go).
func (e Shuffle) validate() error {
	for _, z := range e.Zones {
		if !shuffleableZone(z) {
			return fmt.Errorf("Shuffle: zone %d cannot be shuffled into the deck", z)
		}
	}
	return nil
}

// Text renders the effect: the whole-zone fold "shuffle your hand and discard pile
// into your deck" for Zones, or the bare "shuffle your deck".
func (e Shuffle) Text() string {
	if len(e.Zones) == 0 {
		return "shuffle your deck"
	}
	nouns := make([]string, len(e.Zones))
	for i, z := range e.Zones {
		nouns[i] = z.noun()
	}
	return "shuffle your " + strings.Join(nouns, " and ") + " into your deck"
}

// Resolve folds the chosen zones into the deck (if any) and shuffles: a whole-zone
// fold shuffles once after moving the piles, and a bare shuffle records the
// shuffle itself.
func (e Shuffle) Resolve(ctx *EffectContext) {
	switch {
	case len(e.Zones) > 0:
		ctx.Resolver.ShuffleZonesIntoDeck(ctx.Controller, e.Zones)
	default:
		ctx.Resolver.Shuffle(ctx.Controller)
		ctx.Resolver.Record(DeckShuffled{Player: ctx.Controller})
	}
}

// SwapDeckAndDiscard exchanges the controller's deck with their discard pile and
// shuffles the new deck — Reverse Time turns a spent deck back into a fresh one.
// It differs from Shuffle{Zones: []Zone{Discard}} in that the old deck goes away
// into the discard pile rather than surviving underneath it.
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
