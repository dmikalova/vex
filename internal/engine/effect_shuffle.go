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

// ShuffleChosenCreaturesFromDiscard shuffles any number of creatures the
// controller chooses from their discard pile back into their deck, one grouped
// shuffle (Not Finished with You). An optional House filter narrows the eligible
// creatures.
type ShuffleChosenCreaturesFromDiscard struct {
	House House
}

// Text renders the effect, e.g. "shuffle any number of creatures from your
// discard pile into your deck".
func (e ShuffleChosenCreaturesFromDiscard) Text() string {
	noun := "creatures"
	if e.House != HouseNone {
		noun = e.House.String() + " " + noun
	}
	return "shuffle any number of " + noun + " from your discard pile into your deck"
}

// Resolve lets the controller pick discard-pile creatures one at a time until they
// decline or none remain, shuffling each into their deck as one grouped batch.
func (e ShuffleChosenCreaturesFromDiscard) Resolve(ctx *EffectContext) {
	ctx.Resolver.BeginShuffleBatch()
	picked := map[LocalID]bool{}
	for {
		cands := e.candidates(ctx, picked)
		if len(cands) == 0 {
			break
		}
		id, ok := ctx.ChooseCardOptional("Choose a creature to shuffle into your deck", cands)
		if !ok {
			break
		}
		picked[id] = true
		ctx.Resolver.ShuffleFromDiscardIntoDeck(id)
	}
	ctx.Resolver.EndShuffleBatch(ctx.Source)
}

// candidates lists the controller's not-yet-chosen discard-pile creatures that
// match the optional house filter.
func (e ShuffleChosenCreaturesFromDiscard) candidates(
	ctx *EffectContext,
	picked map[LocalID]bool,
) []LocalID {
	var out []LocalID
	for _, id := range ctx.Resolver.Discard(ctx.Controller) {
		if picked[id] || !ctx.Resolver.IsCreature(id) {
			continue
		}
		if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
			continue
		}
		out = append(out, id)
	}
	return out
}

// ShuffleCardsFromDiscard shuffles a fixed number of cards the controller chooses
// from their discard pile back into their deck, one grouped shuffle. The Count
// says how many — Shard of Life shuffles one card for each friendly Shard. Fewer
// are shuffled when the discard pile holds fewer than the count.
type ShuffleCardsFromDiscard struct {
	Count Count
}

// validate requires a count.
func (e ShuffleCardsFromDiscard) validate() error {
	if e.Count == nil {
		return fmt.Errorf("ShuffleCardsFromDiscard: Count must be set")
	}
	return nil
}

// Text renders the effect, e.g. "shuffle a card from your discard pile into your
// deck for each friendly Shard".
func (e ShuffleCardsFromDiscard) Text() string {
	return forEach(e.Count, "shuffle a card from your discard pile into your deck")
}

// Resolve lets the controller pick that many discard-pile cards one at a time,
// shuffling each into their deck as one grouped batch.
func (e ShuffleCardsFromDiscard) Resolve(ctx *EffectContext) {
	n := e.Count.Value(ctx)
	if n <= 0 {
		return
	}
	ctx.Resolver.BeginShuffleBatch()
	picked := map[LocalID]bool{}
	for i := 0; i < n; i++ {
		var cands []LocalID
		for _, id := range ctx.Resolver.Discard(ctx.Controller) {
			if !picked[id] {
				cands = append(cands, id)
			}
		}
		if len(cands) == 0 {
			break
		}
		id, _ := ctx.ChooseCard("Choose a card to shuffle into your deck", cands)
		picked[id] = true
		ctx.Resolver.ShuffleFromDiscardIntoDeck(id)
	}
	ctx.Resolver.EndShuffleBatch(ctx.Source)
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
