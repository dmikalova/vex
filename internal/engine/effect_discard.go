package engine

import (
	"fmt"
	"slices"
)

// PutCard moves the controller's own cards from one or more source piles
// to a destination — their hand or the top of their deck — with a Selection
// deciding which cards and how: the controller chooses one (Chosen, restrictable
// by type, trait, name, or an Or disjunction), or every matching card is taken
// with no choice (Each). This is how cards recur, e.g. "Put a creature from your
// discard pile on top of your deck." Zones and the destination are both required.
type PutCard struct {
	// Zones names the source piles the cards are taken from. Naming more than one
	// combines them into a single pool — Faygin reaches an Urchin in play or in the
	// discard pile with one choice.
	Zones []Zone
	// Selection decides which cards move and how they are picked; it must be set.
	// Chief Engineer Walls returns a chosen upgrade or Robot card, Ortannu the
	// Chained returns each copy of Ortannu's Binding by name.
	Selection Selection
	// Destination is where the cards go: ToHand or ToTopOfDeck.
	Destination Destination
	// Bind leaves the last card moved in context (ctx.It) so a following effect can
	// act on it — Resurgence returns a creature and, if it is a Mutant, returns
	// another. A selection that moves nothing binds nothing.
	Bind bool
}

// mover is the shared cross-zone plumbing: gather from Zones, move to Destination.
func (e PutCard) mover(ctx *EffectContext) crossZoneMover {
	return crossZoneMover{
		Player:  ctx.Controller,
		Dest:    e.Destination,
		Sources: e.Zones,
	}
}

// sourcePhrase renders where the cards come from, e.g. "from your discard pile".
func (e PutCard) sourcePhrase() string {
	return "from " + whoseZones(Controller, e.Zones)
}

// destPhrase renders where the card goes, e.g. "into your hand".
func (e PutCard) destPhrase() string {
	if e.Destination == ToTopOfDeck {
		return "on top of your deck"
	}
	return "into your hand"
}

// validate rejects an unset selection, an empty zone list, or a destination this
// effect cannot move a card to; only the hand and the top of the deck are
// supported.
func (e PutCard) validate() error {
	if e.Selection == nil {
		return fmt.Errorf("PutCard: selection must be set")
	}
	if len(e.Zones) == 0 {
		return fmt.Errorf("PutCard: Zones must name at least one source")
	}
	if e.Destination != ToHand && e.Destination != ToTopOfDeck {
		return fmt.Errorf("PutCard: unsupported destination %d", e.Destination.zone)
	}
	return nil
}

// Text renders the effect, e.g. "put a creature from your discard pile into your
// hand" or "put each creature of the chosen house from your discard pile into your
// hand".
func (e PutCard) Text() string {
	return "put " + e.Selection.object() + " " + e.sourcePhrase() + " " + e.destPhrase()
}

// listHead, listNoun, and listTail make a type-only recur foldable with its
// neighbours in a Sequence: "put a tactic from your discard pile into your hand"
// and "put an artifact …" fold to "put a tactic, artifact, creature, and upgrade
// from your discard pile into your hand" (Look What I Found!). listNoun is "" for
// any selection that is not a bare card type, so only plain type recurs fold.
func (e PutCard) listHead() string { return "put" }

func (e PutCard) listNoun() string {
	if s, ok := e.Selection.(Chosen); ok && s.plainType() {
		return typeWord(s.Type)
	}
	return ""
}

func (e PutCard) listTail() string {
	return e.sourcePhrase() + " " + e.destPhrase()
}

// vacuous reports that no card in the source piles is a candidate, so a "you may"
// wrapping this effect asks nothing (Chief Engineer Walls prompts only when an
// upgrade or Robot card is actually in the discard).
func (e PutCard) vacuous(ctx *EffectContext) bool {
	return len(e.Selection.candidates(ctx, e.pool(ctx))) == 0
}

// pool gathers every card in the source piles.
func (e PutCard) pool(ctx *EffectContext) []LocalID {
	return e.mover(ctx).gather(ctx, func(LocalID) bool { return true })
}

// Resolve moves the selected cards from the source piles to the destination — a
// Chosen picks one (nothing happens with no candidate), an Each takes every
// match. Each move is tallied for a following
// ProducedThisWay{Tally: TallyCardsReturned}.
func (e PutCard) Resolve(ctx *EffectContext) {
	mover := e.mover(ctx)
	for _, id := range e.Selection.pick(ctx, e.pool(ctx)) {
		mover.move(ctx, id)
		ctx.Produced.Returned++
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
// picked card from either pile (Munchling, Novu Dynamo). Quantity discards that
// many (Old Yurk discards 2), and AnyNumber instead lets the controller discard as
// many matching cards as they like, declining when done (Helmsman Spears). Every
// card discarded this way is recorded on the context so a following
// ForEachDiscarded can act once per card, and it gates a Then (Feeding
// Pit only gains Æmber if a creature was discarded).
type DiscardCard struct {
	// Player whose zone the cards are discarded from.
	Player Player
	// Zones names the source piles the cards are discarded from, each Hand or
	// Archives. Naming both combines them into one pool the controller picks from.
	Zones []Zone
	// Selection decides how each card is picked; it must be set.
	Selection Selection
	// Quantity is how many cards to discard; the zero value discards one.
	Quantity Quantity
	// Bind leaves the last card discarded in context (ctx.It) so a following effect
	// can act on it — Ambassador Liu discards a card and rewards by its house. A
	// discard that moves nothing binds nothing.
	Bind bool
}

// validate rejects a DiscardCard whose player or selection was left unset, an
// empty zone list, or a source zone other than the hand or archives.
func (e DiscardCard) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardCard")
	}
	if e.Selection == nil {
		return fmt.Errorf("DiscardCard: selection must be set")
	}
	if len(e.Zones) == 0 {
		return fmt.Errorf("DiscardCard: at least one zone must be set")
	}
	for _, z := range e.Zones {
		if z != Hand && z != Archives {
			return fmt.Errorf("DiscardCard: zone must be Hand or Archives")
		}
	}
	return quantityValidate(e.Quantity)
}

// object renders the count-bearing noun phrase the discard acts on, e.g. "a card",
// "2 random cards", or "each creature of the chosen house".
func (e DiscardCard) object() string {
	return quantityObject(e.Quantity, e.Selection.noun(), e.Selection.object())
}

// Text renders the effect, naming the source zone explicitly (rule 17). The voice
// is the shared pile-verb template: a blind pick reads "your opponent discards …",
// a controller-directed one "discard … from your opponent's hand" (Deep Probe).
func (e DiscardCard) Text() string {
	return pileVerbText(e.Selection, e.Player, e.Zones, pileVerb{"discard", "discards", e.object()})
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
	limit, bounded := quantityPicks(e.Quantity, ctx)
	for i := 0; !bounded || i < limit; i++ {
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
