// Package deckgen is Vactrol's procedural deck generator. It turns a Set (a pool
// of cards plus tuning) and a seed into a Deck of three House pods of twelve
// Slots each, deterministically. It is a pure function of its inputs and depends
// only on the engine's data types — never the Game runtime, the card facade, or
// provenance (see ADR 0003). The design and vocabulary live in
// docs/deck-generation.md and CONTEXT.md.
//
// What is built today is the core spine: weighted House selection, per-slot
// rarity rolls with Maverick/Special overlays, drawing from the pool, the
// duplicate-pull draw modifier, and per-slot Materialize (identity plus maverick
// rehousing). Legacy pools, connections, the enhancement/distortion finishing
// pass, templates, and scoring are documented seams that stay inert until the
// data and engine support them (a second set, generation profiles, distortions).
package deckgen

import "github.com/dmikalova/vactrol/internal/engine"

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
