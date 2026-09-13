package engine

import "fmt"

// ShuffleFromDiscard shuffles the cards a Selection picks from the controller's
// discard pile back into their deck, as one grouped shuffle. The Selection is the
// one axis that varies (ADR 0031): Each takes every match with no choice (Low Dawn
// returns each Untamed creature), an Optional Chosen with AnyNumber lets the
// controller shuffle back any number of a kind one at a time until they decline
// (Not Finished with You), and a plain Chosen with a Count shuffles that many
// chosen cards (Shard of Life shuffles one for each friendly Shard). Fewer move
// when the discard pile holds fewer than asked. A by-name pick stays its own node,
// ShuffleNamedFromDiscardIntoDeck, because its printed article ("a Subtle Chain")
// belongs to the verb rather than to the Selection (a Named renders bare).
type ShuffleFromDiscard struct {
	// Selection picks which discard-pile cards move; it must be set.
	Selection Selection
	// AnyNumber lets the controller shuffle back as many matching cards as they
	// like, one at a time, until they decline — "any number of". Pair it with an
	// Optional Chosen so the pick is declinable. Mutually exclusive with Count.
	AnyNumber bool
	// Count says how many cards to shuffle when the number is fixed (Shard of Life
	// scales it by the friendly Shards in play). The zero value shuffles one pick.
	Count Count
}

// validate requires a Selection and rejects pairing AnyNumber with a Count (the
// two count modes are mutually exclusive).
func (e ShuffleFromDiscard) validate() error {
	if e.Selection == nil {
		return fmt.Errorf("ShuffleFromDiscard: selection must be set")
	}
	if e.AnyNumber && e.Count != nil {
		return fmt.Errorf("ShuffleFromDiscard: AnyNumber and Count are exclusive")
	}
	return nil
}

// count is the fixed number of picks: the Count's value, or one when unset.
func (e ShuffleFromDiscard) count(ctx *EffectContext) int {
	if e.Count == nil {
		return 1
	}
	return e.Count.Value(ctx)
}

// object renders the noun phrase the shuffle acts on, e.g. "each Untamed creature",
// "any number of creatures", or the singular "a card" a Count wraps with forEach.
func (e ShuffleFromDiscard) object() string {
	if e.AnyNumber {
		return "any number of " + e.Selection.noun() + "s"
	}
	return e.Selection.object()
}

// Text renders the effect, e.g. "shuffle each Untamed creature from your discard
// pile into your deck". A Count leads the sentence with its per-clause.
func (e ShuffleFromDiscard) Text() string {
	base := "shuffle " + e.object() + " from your discard pile into your deck"
	if e.Count != nil {
		return forEach(e.Count, base)
	}
	return base
}

// Resolve shuffles the selected discard-pile cards into the deck as one grouped
// batch, so the several cards narrate one source-attributed line. It re-reads the
// discard pile each pass — a shuffled card has already left it — and stops when a
// pick comes back empty (the pile ran out or an Optional choice was declined).
func (e ShuffleFromDiscard) Resolve(ctx *EffectContext) {
	ctx.Resolver.BeginShuffleBatch()
	for i := 0; e.AnyNumber || i < e.count(ctx); i++ {
		ids := e.Selection.pick(ctx, ctx.Resolver.Discard(ctx.Controller))
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			ctx.Resolver.ShuffleFromDiscardIntoDeck(id)
		}
	}
	ctx.Resolver.EndShuffleBatch()
}
