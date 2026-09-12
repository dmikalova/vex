package deckgen

import (
	"math/rand"

	"github.com/dmikalova/vactrol/internal/engine"
)

// GenerationProfile is the deck-building-only metadata a card carries, kept out
// of the pure engine definition (a facade sidecar, like Provenance). Its zero
// value is an ordinary card. Enhancement sources and multiples are future
// fields; today the flags and the cluster below are read.
type GenerationProfile struct {
	// OneCopyPerDeck bars a second copy in the same deck: once placed, draws and
	// duplicate-pulls skip the card.
	OneCopyPerDeck bool
	// Houseless marks a Special card with no House until it fills a Slot, when it
	// is stamped with that Slot's House.
	Houseless bool
	// RarityWeight scales how often deck generation draws this card among its
	// house+rarity peers, relative to the default weight of 1; a value of 0 (or
	// less) means the default. Five Master-of-N variants at 0.2 draft as often as
	// one ordinary Rare card.
	RarityWeight float64
	// Cluster marks the card a member of a card family placed by strategy, and
	// carries that family's strategy and trigger (ADR 0036). The Shards' cluster
	// (OnePerHouse) places one Shard in each of the deck's Houses whenever any
	// Shard is drawn; the zero value belongs to no cluster.
	Cluster ClusterMembership
}

// SlotContext is what a Materializer needs to produce a concrete card for a Slot.
// House is the pod's House — the card's final House, so a Maverick is rehoused to
// it and a self-house reference binds to it. DeckHouses is the deck's three pod
// Houses (HouseNone for an unfilled pod), so a template can bind a partner house —
// one of the deck's other Houses, as an Ambassador or Plant does (ADR 0036).
type SlotContext struct {
	House      engine.House
	Rarity     engine.Rarity
	Maverick   bool
	Legacy     bool
	Special    bool
	DeckHouses [PodCount]engine.House
}

// Materializer turns a pool entry into a concrete, engine-ready card at
// generation time (see ADR 0004). Concrete cards use the identity materializer (a
// nil Materializer on a Card); a template binds its parameters, name, and
// self-house references here — an Ambassador or Plant keys on the pod's partner
// House, a bane on three random Houses. The returned definition must be flat and
// pointerless, exactly what the engine consumes.
type Materializer interface {
	Materialize(ctx SlotContext, r *rand.Rand) engine.CardDefinition
}
