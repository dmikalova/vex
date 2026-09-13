package engine

import (
	"fmt"
)

// PutFromDiscard moves a card the controller chooses from their own discard pile
// to a destination — their hand or the top of their deck. Type restricts the
// choice to cards of that type; the zero value allows any card. With All it moves
// every matching card instead of one chosen card (Arise! returning each creature
// of a house). This is how cards recur from the discard pile, e.g. "Put a creature
// from your discard pile on top of your deck." The destination is required.
type PutFromDiscard struct {
	// Match restricts the choice to cards the predicate admits; the zero value
	// admits any card. Chief Engineer Walls returns an upgrade or Robot card,
	// Ortannu the Chained returns each copy of Ortannu's Binding by name.
	Match Match
	// Destination is where the card goes: ToHand or ToTopOfDeck.
	Destination Destination
	// All moves every matching card instead of one chosen card (Arise! returning
	// each creature of a house).
	All bool
	// OfChosenHouse limits the matching cards to the house an enclosing
	// ChooseHouseThen picked. It applies only with All.
	OfChosenHouse bool
}

// noun renders the kind of card the effect moves, delegating to the match.
func (e PutFromDiscard) noun() string {
	return e.Match.noun()
}

// destPhrase renders where the card goes, e.g. "into your hand".
func (e PutFromDiscard) destPhrase() string {
	if e.Destination == ToTopOfDeck {
		return "on top of your deck"
	}
	return "into your hand"
}

// validate rejects a destination this effect cannot move a card to; only the hand
// and the top of the deck are supported, and the destination must be named.
func (e PutFromDiscard) validate() error {
	if e.Destination != ToHand && e.Destination != ToTopOfDeck {
		return fmt.Errorf("PutFromDiscard: unsupported destination %d", e.Destination.zone)
	}
	return nil
}

// Text renders the effect, e.g. "put a card from your discard pile into your hand"
// or "put each creature of the chosen house from your discard pile into your hand".
func (e PutFromDiscard) Text() string {
	if e.All {
		what := "each " + e.noun()
		if e.OfChosenHouse {
			what += " of the chosen house"
		}
		return "put " + what + " from your discard pile " + e.destPhrase()
	}
	return "put " + indefinite(e.noun()) + " from your discard pile " + e.destPhrase()
}

// moveTo moves one card from the discard pile to the destination and tallies it
// for a following ProducedThisWay{Tally: TallyCardsReturned}.
func (e PutFromDiscard) moveTo(ctx *EffectContext, id LocalID) {
	if e.Destination == ToTopOfDeck {
		ctx.Resolver.MoveFromDiscardToTopOfDeck(id)
	} else {
		ctx.Resolver.PutFromDiscardIntoHand(id)
	}
	ctx.Produced.Returned++
}

// admits reports whether a discard-pile card passes the Match filter (the
// OfChosenHouse filter is applied separately, only with All).
func (e PutFromDiscard) admits(ctx *EffectContext, id LocalID) bool {
	return e.Match.admits(ctx.Resolver, id)
}

// matches reports whether a discard-pile card is a candidate this effect could
// move — the admits filters plus the OfChosenHouse restriction (which applies only
// with All).
func (e PutFromDiscard) matches(ctx *EffectContext, id LocalID) bool {
	return e.admits(ctx, id) &&
		(!e.All || !e.OfChosenHouse || ctx.Resolver.House(id) == ctx.ChosenHouse)
}

// vacuous reports that no card in the controller's discard pile matches, so a
// "you may" wrapping this effect asks nothing (Chief Engineer Walls prompts only
// when an upgrade or Robot card is actually in the discard).
func (e PutFromDiscard) vacuous(ctx *EffectContext) bool {
	return len(discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		return e.matches(ctx, id)
	})) == 0
}

// Resolve moves a card from the controller's discard pile to the destination. With
// All it moves every matching card; otherwise the controller chooses one, and
// nothing happens if there is no candidate or the choice is declined.
func (e PutFromDiscard) Resolve(ctx *EffectContext) {
	if e.All {
		for _, id := range discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
			return e.matches(ctx, id)
		}) {
			e.moveTo(ctx, id)
		}
		return
	}
	candidates := discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		return e.matches(ctx, id)
	})
	id, ok := ctx.ChooseCreature("Choose a "+e.noun()+" from your discard pile", candidates)
	if !ok {
		return
	}
	e.moveTo(ctx, id)
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
			if fromArchives {
				ctx.Resolver.DiscardCardFromArchives(owner, id)
			} else {
				ctx.Resolver.DiscardCardFromHand(owner, id)
			}
			ctx.Produced.Discarded = append(ctx.Produced.Discarded, id)
			moved = true
		}
	}
	return moved
}
