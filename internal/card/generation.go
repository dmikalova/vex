package card

import (
	"math/rand"

	"github.com/dmikalova/vactrol/internal/deckgen"
)

// Deck-generation metadata re-exported for authoring. These attach to a card the
// same way Provenance does — recorded in the registry, not on the engine
// definition — and are read by the deck generator (see internal/deckgen and
// ADR 0004).
type (
	// Materializer turns a card template into a concrete card at generation time.
	Materializer = deckgen.Materializer
	// SlotContext is what a Materializer is given: the pod's house and the slot's
	// rolled rarity and provenance flags.
	SlotContext = deckgen.SlotContext
	// Connection is the set of connected cards a puller card brings into its pod;
	// build one with card.Connects.
	Connection = deckgen.Connection
)

// MaterializeFunc adapts a plain function to a Materializer, so a template can be
// written inline: card.Template(func(ctx card.SlotContext, r *rand.Rand) card.Definition { ... }).
type MaterializeFunc func(SlotContext, *rand.Rand) Definition

// Materialize satisfies the Materializer interface.
func (f MaterializeFunc) Materialize(ctx SlotContext, r *rand.Rand) Definition { return f(ctx, r) }

// Template marks a card as a template: a single pool entry whose concrete cards
// deck generation produces from f. The card.New face (house, type, rarity) buckets
// the template in the pool; f builds the playable card.
func Template(f MaterializeFunc) Option { return func(b *builder) { b.materializer = f } }

// OneCopyPerDeck bars deck generation from placing more than one copy of the card.
func OneCopyPerDeck() Option { return func(b *builder) { b.profile.OneCopyPerDeck = true } }

// Connects marks the card as a connection puller: when it is placed in a pod,
// deck generation pulls one copy of each given card into that pod. The connected
// cards are named by their own definition symbols (e.g. card.Connects(
// HelpFromFutureSelf)), so a connection to a card that does not exist is a
// compile error, never a silently dropped link. Each connected card must be
// authored with card.Rarity.Connected so it never rolls on its own.
func Connects(cards ...Definition) Option {
	names := make([]string, len(cards))
	for i, c := range cards {
		names[i] = c.Name
	}
	return func(b *builder) { b.profile.Connection = deckgen.Connection{Cards: names} }
}
