package engine

import (
	"fmt"
	"slices"
	"strings"
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
	// Bind leaves the last card moved in context (ctx.It) so a following effect can
	// act on it — Resurgence returns a creature and, if it is a Mutant, returns
	// another. A selection that moves nothing binds nothing.
	Bind bool
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

// listHead, listNoun, and listTail make a type-only recur foldable with its
// neighbours in a Sequence: "put a tactic from your discard pile into your hand"
// and "put an artifact …" fold to "put a tactic, artifact, creature, and upgrade
// from your discard pile into your hand" (Look What I Found!). listNoun is "" for
// any selection that is not a bare card type, so only plain type recurs fold.
func (e PutFromDiscard) listHead() string { return "put" }

func (e PutFromDiscard) listNoun() string {
	if s, ok := e.Selection.(Chosen); ok && s.plainType() {
		return typeWord(s.Type)
	}
	return ""
}

func (e PutFromDiscard) listTail() string {
	return "from your discard pile " + e.destPhrase()
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
		if e.Bind {
			ctx.It, ctx.HasIt = id, true
		}
	}
}

// DiscardCard discards cards from a player's hand or archives, with a Selection
// deciding how each card is picked — the controller chooses one (Chosen,
// restrictable to a type), or a uniformly random card leaves a hidden zone
// (Random, for Mind Barb and Tantadlin). Zones names the source piles the discard
// draws from, each Hand or Archives; naming both lets the controller discard the
// picked card from either pile (Munchling, Novu Dynamo). Amount discards that many
// (Old Yurk discards 2); the zero value discards one. AnyNumber instead lets the
// controller discard as many matching cards as they like, declining when done
// (Helmsman Spears). Every card discarded this way is recorded on the context so a
// following ForEachDiscarded can act once per card, and it gates a Then (Feeding
// Pit only gains Æmber if a creature was discarded).
type DiscardCard struct {
	// Player whose zone the cards are discarded from.
	Player Player
	// Zones names the source piles the cards are discarded from, each Hand or
	// Archives. Naming both combines them into one pool the controller picks from.
	Zones []Zone
	// Selection decides how each card is picked; it must be set.
	Selection Selection
	// Amount is how many cards to discard; the zero value counts as one.
	Amount int
	// AnyNumber lets the controller discard as many cards as they like instead of a
	// fixed Amount, declining when done.
	AnyNumber bool
	// Bind leaves the last card discarded in context (ctx.It) so a following effect
	// can act on it — Ambassador Liu discards a card and rewards by its house. A
	// discard that moves nothing binds nothing.
	Bind bool
}

// validate rejects a DiscardCard whose player or selection was left unset, an
// Amount paired with AnyNumber (the two count modes are mutually exclusive), an
// empty zone list, or a source zone other than the hand or archives.
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
	if len(e.Zones) == 0 {
		return fmt.Errorf("DiscardCard: at least one zone must be set")
	}
	for _, z := range e.Zones {
		if z != Hand && z != Archives {
			return fmt.Errorf("DiscardCard: zone must be Hand or Archives")
		}
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

// sourceNoun names the piles the discard draws from, joining them with "or" when
// the controller may pick from either (Munchling's "hand or archives").
func (e DiscardCard) sourceNoun() string {
	nouns := make([]string, len(e.Zones))
	for i, z := range e.Zones {
		nouns[i] = z.noun()
	}
	return strings.Join(nouns, " or ")
}

// Text renders the effect, naming the source zone explicitly (rule 17). The voice
// depends on who discards: a random pick from a player's hidden zone reads "your
// opponent discards …", while a discard the controller directs reads "discard …
// from your opponent's hand" (Deep Probe).
func (e DiscardCard) Text() string {
	obj := e.object()
	zone := e.sourceNoun()
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
	moved := false
	for i := 0; e.AnyNumber || i < e.count(); i++ {
		ids := e.Selection.pick(ctx, e.source(ctx, owner))
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			toDiscard.moveFrom(ctx, e.zoneOf(ctx, owner, id), owner, id)
			recordDiscardedThisWay(ctx, id)
			moved = true
			if e.Bind {
				ctx.It, ctx.HasIt = id, true
			}
		}
	}
	return moved
}

// zoneCards returns the cards in one of the discard's source piles.
func (e DiscardCard) zoneCards(ctx *EffectContext, owner int, z Zone) []LocalID {
	if z == Archives {
		return ctx.Resolver.Archives(owner)
	}
	return ctx.Resolver.Hand(owner)
}

// source returns the cards the selection may pick from: one pile on its own, or —
// when Zones names both — the hand and archives combined (Munchling).
func (e DiscardCard) source(ctx *EffectContext, owner int) []LocalID {
	if len(e.Zones) == 1 {
		return e.zoneCards(ctx, owner, e.Zones[0])
	}
	var combined []LocalID
	for _, z := range e.Zones {
		combined = append(combined, e.zoneCards(ctx, owner, z)...)
	}
	return combined
}

// zoneOf reports which pile a picked card sits in, so a combined hand-or-archives
// discard removes it from the right one. Every pile but the last is checked; a card
// found in none of them must sit in the last, which also covers a single-zone discard.
func (e DiscardCard) zoneOf(ctx *EffectContext, owner int, id LocalID) Zone {
	for _, z := range e.Zones[:len(e.Zones)-1] {
		if slices.Contains(e.zoneCards(ctx, owner, z), id) {
			return z
		}
	}
	return e.Zones[len(e.Zones)-1]
}
