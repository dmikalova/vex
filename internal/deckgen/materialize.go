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
	// placed (Timetraveller pulls Help from Future Self; Troop Call pulls the Niffle
	// Apes it calls). See Connection.
	Connection Connection
}

// Connection is the set of connected cards a puller card brings into its pod.
// Each entry is ensured present at its copy count, overwriting other
// (unprotected) slots. A maverick puller still fires its connection, and its
// partners are rehoused to the pod's House along with it (the printed-house
// rule KeyForge itself uses for a Maverick's connected cards).
type Connection struct {
	// Cards are the connected cards, each pulled at its own count and rate.
	Cards []ConnectedCard
}

// ConnectedCard is one card a connection pulls: how many copies the pod ends up
// holding, and how often the pull happens at all. A guaranteed partner
// (Timetraveller's Help from Future Self) carries Chance 1; a flavourful one
// (Troop Call's Niffle Queen) carries less, and is rolled once per pod.
//
// A connected card need not be rarity Connected: Troop Call guarantees Niffle
// Apes that roll on their own too. Author a card Connected only when it should
// never appear without its puller, since the pool skips those entirely.
type ConnectedCard struct {
	// Name is the connected card's name, as it appears in the set.
	Name string
	// Copies is how many of it the pod ends up holding; at least one.
	Copies int
	// Chance is how often the pull fires, in (0, 1]; 1 is every time.
	Chance float64
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
