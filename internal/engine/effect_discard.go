package engine

import (
	"fmt"
)

// PutFromDiscard moves cards from the controller's own discard pile to a
// destination — their hand or the top of their deck — with a Selection deciding
// which cards and how: the controller chooses one (Chosen, restrictable by type,
// trait, name, or an Or disjunction), or every matching card is taken with no
// choice (Each). This is how cards recur from the discard pile, e.g. "Put a
// creature from your discard pile on top of your deck." The destination is
// required.
type PutFromDiscard struct {
	// Selection decides which discard-pile cards move and how they are picked; it
	// must be set. Chief Engineer Walls returns a chosen upgrade or Robot card,
	// Ortannu the Chained returns each copy of Ortannu's Binding by name.
	Selection Selection
	// Destination is where the cards go: ToHand or ToTopOfDeck.
	Destination Destination
}

// destPhrase renders where the card goes, e.g. "into your hand".
func (e PutFromDiscard) destPhrase() string {
	if e.Destination == ToTopOfDeck {
		return "on top of your deck"
	}
	return "into your hand"
}

// validate rejects an unset selection or a destination this effect cannot move a
// card to; only the hand and the top of the deck are supported.
func (e PutFromDiscard) validate() error {
	if e.Selection == nil {
		return fmt.Errorf("PutFromDiscard: selection must be set")
	}
	if e.Destination != ToHand && e.Destination != ToTopOfDeck {
		return fmt.Errorf("PutFromDiscard: unsupported destination %d", e.Destination.zone)
	}
	return nil
}

// Text renders the effect, e.g. "put a creature from your discard pile into your
// hand" or "put each creature of the chosen house from your discard pile into your
// hand".
func (e PutFromDiscard) Text() string {
	return "put " + e.Selection.object() + " from your discard pile " + e.destPhrase()
}

// moveTo moves one card from the discard pile to the destination and tallies it
// for a following ProducedThisWay{Tally: TallyCardsReturned}.
func (e PutFromDiscard) moveTo(ctx *EffectContext, id LocalID) {
	e.Destination.moveFrom(ctx, Discard, ctx.Controller, id)
	ctx.Produced.Returned++
}

// vacuous reports that no card in the controller's discard pile is a candidate, so
// a "you may" wrapping this effect asks nothing (Chief Engineer Walls prompts only
// when an upgrade or Robot card is actually in the discard).
func (e PutFromDiscard) vacuous(ctx *EffectContext) bool {
	return len(e.Selection.candidates(ctx, ctx.Resolver.Discard(ctx.Controller))) == 0
}

// Resolve moves the selected cards from the controller's discard pile to the
// destination — a Chosen picks one (nothing happens with no candidate), an Each
// takes every match.
func (e PutFromDiscard) Resolve(ctx *EffectContext) {
	for _, id := range e.Selection.pick(ctx, ctx.Resolver.Discard(ctx.Controller)) {
		e.moveTo(ctx, id)
	}
}

// DiscardCard discards cards from a player's hand or archives, with a Selection
// deciding how each card is picked — the controller chooses one (Chosen,
// restrictable to a type), or a uniformly random card leaves a hidden zone
// (Random, for Mind Barb and Tantadlin). Zone names the source, Hand or Archives.
// Amount discards that many (Old Yurk discards 2); the zero value discards one.
// AnyNumber instead lets the controller discard as many matching cards as they
// like, declining when done (Helmsman Spears). Every card discarded this way is
// recorded on the context so a following ForEachDiscarded can act once per card,
// and it gates a Then (Feeding Pit only gains Æmber if a creature was discarded).
type DiscardCard struct {
	// Player whose zone the cards are discarded from.
	Player Player
	// Zone names the source the cards are discarded from: Hand or Archives.
	Zone Zone
	// Selection decides how each card is picked; it must be set.
	Selection Selection
	// Amount is how many cards to discard; the zero value counts as one.
	Amount int
	// AnyNumber lets the controller discard as many cards as they like instead of a
	// fixed Amount, declining when done.
	AnyNumber bool
}

// validate rejects a DiscardCard whose player or selection was left unset, an
// Amount paired with AnyNumber (the two count modes are mutually exclusive), or a
// source zone other than the hand or archives.
func (e DiscardCard) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardCard")
	}
	if e.Selection == nil {
		return fmt.Errorf("DiscardCard: selection must be set")
	}
	if e.AnyNumber && e.Amount != 0 {
		return fmt.Errorf("DiscardCard: AnyNumber and Amount are exclusive")
	}
	if e.Zone != Hand && e.Zone != Archives {
		return fmt.Errorf("DiscardCard: zone must be Hand or Archives")
	}
	return nil
}

// count is Amount with the zero value treated as one.
func (e DiscardCard) count() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// object renders the count-bearing noun phrase the discard acts on, e.g. "a card",
// "2 random cards", or "each creature of the chosen house".
func (e DiscardCard) object() string {
	switch {
	case e.AnyNumber:
		return "any number of " + e.Selection.noun() + "s"
	case e.count() > 1:
		return countNoun(e.count(), e.Selection.noun())
	default:
		return e.Selection.object()
	}
}

// Text renders the effect, naming the source zone explicitly (rule 17). The voice
// depends on who discards: a random pick from a player's hidden zone reads "your
// opponent discards …", while a discard the controller directs reads "discard …
// from your opponent's hand" (Deep Probe).
func (e DiscardCard) Text() string {
	obj := e.object()
	zone := e.Zone.noun()
	switch e.Player {
	case Opponent:
		if selectionOwnerActs(e.Selection) {
			return "your opponent discards " + obj + " from their " + zone
		}
		return "discard " + obj + " from your opponent's " + zone
	case ItsOwner:
		return "its owner discards " + obj + " from their " + zone
	default:
		return "discard " + obj + " from your " + zone
	}
}

// Resolve discards the selected cards from the player's zone, stopping early if
// the source runs out or a choice is declined.
func (e DiscardCard) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate performs the discards and reports whether any card was discarded, so
// DiscardCard can gate a Then. Each discarded card is appended to
// ctx.Produced.Discarded and moved through the zone's discard method — from the
// hand it fires the after-discard reactions; from the hidden archives it does not.
func (e DiscardCard) resolveGate(ctx *EffectContext) bool {
	owner := ctx.PlayerFor(e.Player)
	fromArchives := e.Zone == Archives
	moved := false
	for i := 0; e.AnyNumber || i < e.count(); i++ {
		source := ctx.Resolver.Hand(owner)
		if fromArchives {
			source = ctx.Resolver.Archives(owner)
		}
		ids := e.Selection.pick(ctx, source)
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			toDiscard.moveFrom(ctx, e.Zone, owner, id)
			ctx.Produced.Discarded = append(ctx.Produced.Discarded, id)
			moved = true
		}
	}
	return moved
}
