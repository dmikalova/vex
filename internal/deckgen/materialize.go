package deckgen

import (
	"math/rand"

	"github.com/dmikalova/vactrol/internal/engine"
)

// GenerationProfile is the deck-building-only metadata a card carries, kept out
// of the pure engine definition (a facade sidecar, like Provenance). Its zero
// value is an ordinary card. Enhancement sources and multiples are future
// fields; today the flags and the connection below are read.
type GenerationProfile struct {
	// OneCopyPerDeck bars a second copy in the same deck: once placed, draws and
	// duplicate-pulls skip the card.
	OneCopyPerDeck bool
	// Houseless marks a Special card with no House until it fills a Slot, when it
	// is stamped with that Slot's House.
	Houseless bool
	// Connection names the connected cards this card pulls into its pod when it is
	// placed (Timetraveller pulls Help from Future Self; Horseman of Pestilence
	// pulls the other three Horsemen). See Connection.
	Connection Connection
}

// Connection is the set of connected cards a puller card brings into its pod.
// Each named card is ensured present in the pod exactly once, overwriting other
// (unprotected) slots. The connected cards carry rarity Connected so they never
// roll on their own. Duplicate counts and cross-house (maverick) connections are
// future axes; today a connection ensures one of each named card, in-house only.
type Connection struct {
	// Cards are the connected cards' names, each pulled into the pod once.
	Cards []string
}

// Empty reports whether the connection pulls nothing.
func (c Connection) Empty() bool { return len(c.Cards) == 0 }

// SlotContext is what a Materializer needs to produce a concrete card for a Slot.
// House is the pod's House — the card's final House, so a Maverick is rehoused to
// it and a self-house reference binds to it.
type SlotContext struct {
	House    engine.House
	Rarity   engine.Rarity
	Maverick bool
	Legacy   bool
	Special  bool
}

// Materializer turns a pool entry into a concrete, engine-ready card at
// generation time (see ADR 0004). Concrete cards use the identity materializer (a
// nil Materializer on a Card); templates — a future addition — bind their
// parameters, name, and self-house references here. The returned definition must
// be flat and pointerless, exactly what the engine consumes.
type Materializer interface {
	Materialize(ctx SlotContext, r *rand.Rand) engine.CardDefinition
}
