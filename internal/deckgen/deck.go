// Package deckgen is Vactrol's procedural deck generator. It turns a Set (a pool
// of cards plus tuning) and a seed into a Deck of three House pods of twelve
// Slots each, deterministically. It is a pure function of its inputs and depends
// only on the engine's data types — never the Game runtime, the card facade, or
// provenance (see ADR 0003). The design and vocabulary live in
// docs/deck-generation.md and CONTEXT.md.
//
// What is built today is the core spine: weighted House selection, per-slot
// rarity rolls with Maverick/Special overlays, drawing from the pool, the
// duplicate-pull draw modifier, per-slot Materialize (identity, maverick
// rehousing, and templates that bind a card per pod — Ambassadors, Plants,
// banes), cross-set Legacy pools, cluster placement (ADR 0036), and the deck-wide
// Enhance finishing pass (ADR 0004) that distributes each Enhance source's bonus
// icons onto random cards. The distortion finishing pass and scoring remain
// documented seams that stay inert until the engine supports them.
package deckgen

import (
	"fmt"

	"github.com/dmikalova/vactrol/internal/engine"
)

// A Deck is three House pods of twelve Slots — 36 cards — generated from one Set
// and seed. It is reproducible only within a single version of its Set's pool.
const (
	// PodCount is the number of House pods in a deck.
	PodCount = 3
	// PodSize is the number of Slots in a House pod.
	PodSize = 12
	// DeckSize is the total number of cards in a deck.
	DeckSize = PodCount * PodSize
)

// Deck is the generated result: three House pods plus the inputs that produced
// it, so it can be reproduced.
type Deck struct {
	Set  string
	Seed int64
	Pods [PodCount]HousePod
}

// HousePod is one of a Deck's three Houses together with its twelve Slots.
type HousePod struct {
	House engine.House
	Slots [PodSize]Slot
}

// Slot is one of a House pod's twelve positions: its intrinsic Rarity, its
// provenance flags, and the materialized card that fills it. Card is the final
// playable definition (a Maverick is already rehoused to the pod's House).
type Slot struct {
	Rarity   engine.Rarity
	Maverick bool
	Legacy   bool
	Special  bool
	Card     engine.CardDefinition
}

// Cards returns the deck's 36 card definitions, pod by pod, in order.
func (d Deck) Cards() []engine.CardDefinition {
	out := make([]engine.CardDefinition, 0, DeckSize)
	for _, pod := range d.Pods {
		for _, s := range pod.Slots {
			out = append(out, s.Card)
		}
	}
	return out
}

// Houses returns the deck's three pod Houses, in pod order (sorted by name).
func (d Deck) Houses() []engine.House {
	hs := make([]engine.House, 0, PodCount)
	for _, pod := range d.Pods {
		hs = append(hs, pod.House)
	}
	return hs
}

// validate asserts the finished deck is well-formed after every fill and cluster
// pass (ADR 0036): every card sits in a slot of a pod whose House it shares. The
// fixed [PodCount][PodSize] shape already guarantees the DeckSize count and that
// no pod is over-full; cluster expansion overwrites slots rather than adding them,
// so it cannot change the count either. What it (and maverick rehousing) can break
// is House integrity — forcing a card of the wrong House into a pod — so that is
// what this checks, per filled slot of a real pod. A House-less pod or an empty
// slot is a degenerate-pool artifact (fewer than PodCount Houses, or a pool too
// small to fill a pod) that the cluster pass also skips, so it is passed over here
// too. validate panics rather than return a mis-housed deck, so a cluster or fill
// regression fails loudly at generation time instead of emitting an illegal deck.
func (d Deck) validate() {
	for i := range d.Pods {
		pod := d.Pods[i]
		if pod.House == engine.HouseNone {
			continue
		}
		for j := range pod.Slots {
			card := pod.Slots[j].Card
			if card.Name == "" {
				continue
			}
			if card.House != pod.House {
				panic(fmt.Sprintf(
					"deckgen: pod %d slot %d holds %q of House %s, not the pod's House %s",
					i, j, card.Name, card.House, pod.House,
				))
			}
		}
	}
}
