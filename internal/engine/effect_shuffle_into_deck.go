package engine

import (
	"fmt"
	"slices"
)

// ShuffleIntoDeck shuffles the cards a Selection picks from the controller's own
// zones back into their deck, as one grouped shuffle. Three axes vary (ADR 0031).
// From names the source zones — the discard pile alone for most cards, or a set a
// single pick ranges over (Song of Spring reaches hand, discard pile, and
// battleline at once). Selection picks within them: Each takes every match with no
// choice (Low Dawn returns each Untamed creature), an Optional Chosen lets the
// controller stop early, and a Named takes the one card of a printed name (Chain
// Gang returns Subtle Chain). Quantity says how many: AnyNumber shuffles back as
// many as the controller likes (Not Finished with You), Scaled shuffles one per
// board count (Shard of Life, one for each friendly Shard). Fewer move when the
// zones hold fewer than asked.
type ShuffleIntoDeck struct {
	// Player names whose zones the cards are taken from; it must be set.
	Player Player
	// From names the zones the cards are drawn from, rendered in this order; it
	// must list at least one of Hand, Discard, or InPlay.
	From []Zone
	// Selection picks which of those zones' cards move; it must be set.
	Selection Selection
	// Quantity says how many cards move; the zero value shuffles one pick.
	Quantity Quantity
}

// validate requires a Selection and at least one supported source zone.
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
	return quantityValidate(e.Quantity)
}

// side renders the scope an in-play source has to spell out. The mover reaches
// only the Player's own cards, but the zone phrase "play" names no side, so the
// object carries "friendly" (or "enemy") instead. Empty when no source is play.
func (e ShuffleIntoDeck) side() string {
	if !slices.Contains(e.From, InPlay) {
		return ""
	}
	return side(e.Player)
}

// object renders the noun phrase the shuffle acts on, e.g. "each friendly card",
// "any number of creatures", or the singular "a card" a Scaled quantity leads with.
func (e ShuffleIntoDeck) object() string {
	side := e.side()
	return quantityObject(
		e.Quantity,
		qualifyNoun(side, e.Selection.noun()),
		selectionObject(e.Selection, side))
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
// pile into your deck". A Scaled quantity leads the sentence with its per-clause.
// The deck is always the same player's, so it takes the same possessive as the
// sources.
func (e ShuffleIntoDeck) Text() string {
	base := "shuffle " + e.object() + " from " + e.fromText() +
		" into " + possessive(e.Player) + " deck"
	if per := quantityLeadIn(e.Quantity); per != nil {
		return forEach(per, base)
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
	limit, bounded := quantityPicks(e.Quantity, ctx)
	for i := 0; !bounded || i < limit; i++ {
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
