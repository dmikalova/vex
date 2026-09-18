package engine

import "fmt"

// ShuffleIntoDeck shuffles the cards a Selection picks from the controller's own
// zones back into their deck, as one grouped shuffle. Two axes vary (ADR 0031).
// From names the source zones — the discard pile alone for most cards, or a set a
// single pick ranges over (Song of Spring reaches hand, discard pile, and
// battleline at once). Selection picks within them: Each takes every match with no
// choice (Low Dawn returns each Untamed creature), an Optional Chosen with
// AnyNumber lets the controller shuffle back any number one at a time until they
// decline (Not Finished with You), a plain Chosen with a Count shuffles that many
// chosen cards (Shard of Life shuffles one for each friendly Shard), and a Named
// takes the one card of a printed name (Chain Gang returns Subtle Chain). Fewer
// move when the zones hold fewer than asked.
type ShuffleIntoDeck struct {
	// Player names whose zones the cards are taken from; it must be set.
	Player Player
	// From names the zones the cards are drawn from, rendered in this order; it
	// must list at least one of Hand, Discard, or InPlay.
	From []Zone
	// Selection picks which of those zones' cards move; it must be set.
	Selection Selection
	// AnyNumber lets the controller shuffle back as many matching cards as they
	// like, one at a time, until they decline — "any number of". Pair it with an
	// Optional Chosen so the pick is declinable. Mutually exclusive with Count.
	AnyNumber bool
	// Count says how many cards to shuffle when the number is fixed (Shard of Life
	// scales it by the friendly Shards in play). The zero value shuffles one pick.
	Count Count
}

// validate requires a Selection and at least one supported source zone, and
// rejects pairing AnyNumber with a Count (the two count modes are exclusive).
func (e ShuffleIntoDeck) validate() error {
	if e.Player == playerUnset {
		return fmt.Errorf("ShuffleIntoDeck: player must be set")
	}
	if e.Selection == nil {
		return fmt.Errorf("ShuffleIntoDeck: selection must be set")
	}
	if len(e.From) == 0 {
		return fmt.Errorf("ShuffleIntoDeck: From must name at least one zone")
	}
	for _, z := range e.From {
		if z != Hand && z != Discard && z != InPlay {
			return fmt.Errorf(
				"ShuffleIntoDeck: From must be Hand, Discard, or InPlay")
		}
	}
	if e.AnyNumber && e.Count != nil {
		return fmt.Errorf("ShuffleIntoDeck: AnyNumber and Count are exclusive")
	}
	return nil
}

// count is the fixed number of picks: the Count's value, or one when unset.
func (e ShuffleIntoDeck) count(ctx *EffectContext) int {
	if e.Count == nil {
		return 1
	}
	return e.Count.Value(ctx)
}

// object renders the noun phrase the shuffle acts on, e.g. "each Untamed creature",
// "any number of creatures", or the singular "a card" a Count wraps with forEach.
func (e ShuffleIntoDeck) object() string {
	if e.AnyNumber {
		return "any number of " + e.Selection.noun() + "s"
	}
	return e.Selection.object()
}

// fromText renders the source zones as an alternation of fully formed phrases,
// e.g. "your discard pile" or "your hand, your discard pile, or in play". Each
// phrase carries its own determiner because in play takes none.
func (e ShuffleIntoDeck) fromText() string {
	phrases := make([]string, len(e.From))
	for i, z := range e.From {
		phrases[i] = whoseZone(e.Player, z)
	}
	return joinOr(phrases)
}

// Text renders the effect, e.g. "shuffle each Untamed creature from your discard
// pile into your deck". A Count leads the sentence with its per-clause. The deck
// is always the same player's, so it takes the same possessive as the sources.
func (e ShuffleIntoDeck) Text() string {
	base := "shuffle " + e.object() + " from " + e.fromText() +
		" into " + possessive(e.Player) + " deck"
	if e.Count != nil {
		return forEach(e.Count, base)
	}
	return base
}

// Resolve shuffles the selected cards into the deck as one grouped batch, so the
// several cards narrate one source-attributed line. It re-reads the source zones
// each pass — a shuffled card has already left them — and stops when a pick comes
// back empty (the zones ran out or an Optional choice was declined).
//
// Each moved card is tallied on ctx.Produced.Moved under its *owner*, not under
// the player whose zones it came from: a card returns to its owner's deck, so a
// controlled enemy creature shuffled home refills the opponent's deck and must
// not feed a following Draw{Per: CardsShuffledIntoDeck} (Timequake). Pinned by
// TestShuffleIntoDeckTalliesByOwner.
func (e ShuffleIntoDeck) Resolve(ctx *EffectContext) {
	mover := crossZoneMover{
		Player:  ctx.PlayerFor(e.Player),
		Dest:    ToDeckShuffled,
		Sources: e.From,
	}
	ctx.Resolver.BeginShuffleBatch()
	for i := 0; e.AnyNumber || i < e.count(ctx); i++ {
		ids := e.Selection.pick(ctx, mover.gather(ctx, func(LocalID) bool { return true }))
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			owner := ctx.Resolver.Owner(id)
			mover.move(ctx, id)
			ctx.Produced.Moved[owner]++
		}
	}
	ctx.Resolver.EndShuffleBatch()
}
