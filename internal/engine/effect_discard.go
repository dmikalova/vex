package engine

import (
	"fmt"
	"strings"
)

// PutFromDiscard moves a card the controller chooses from their own discard pile
// to a destination — their hand or the top of their deck. Type restricts the
// choice to cards of that type; the zero value allows any card. With All it moves
// every matching card instead of one chosen card (Arise! returning each creature
// of a house). This is how cards recur from the discard pile, e.g. "Put a creature
// from your discard pile on top of your deck." The destination is required.
type PutFromDiscard struct {
	// Type restricts the choice to cards of that type; the zero value (an unset
	// CardType) allows any card.
	Type CardType
	// Destination is where the card goes: ToHand or ToTopOfDeck.
	Destination Destination
	// All moves every matching card instead of one chosen card (Arise! returning
	// each creature of a house).
	All bool
	// OfChosenHouse limits the matching cards to the house an enclosing
	// ChooseHouseThen picked. It applies only with All.
	OfChosenHouse bool
}

// noun renders the kind of card the effect moves — the lowercased card type when
// Type is set (e.g. "creature"), otherwise the generic "card".
func (e PutFromDiscard) noun() string {
	if e.Type != "" {
		return strings.ToLower(string(e.Type))
	}
	return "card"
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
		return fmt.Errorf("PutFromDiscard: unsupported destination %d", e.Destination)
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
	return "put a " + e.noun() + " from your discard pile " + e.destPhrase()
}

// moveTo moves one card from the discard pile to the destination.
func (e PutFromDiscard) moveTo(ctx *EffectContext, id LocalID) {
	if e.Destination == ToTopOfDeck {
		ctx.Resolver.MoveFromDiscardToTopOfDeck(id)
	} else {
		ctx.Resolver.PutFromDiscardIntoHand(id)
	}
}

// Resolve moves a card from the controller's discard pile to the destination. With
// All it moves every matching card; otherwise the controller chooses one, and
// nothing happens if there is no candidate or the choice is declined.
func (e PutFromDiscard) Resolve(ctx *EffectContext) {
	if e.All {
		for _, id := range ctx.Resolver.Discard(ctx.Controller) {
			if e.Type != "" && ctx.Resolver.TypeOf(id) != e.Type {
				continue
			}
			if e.OfChosenHouse && ctx.Resolver.House(id) != ctx.ChosenHouse {
				continue
			}
			e.moveTo(ctx, id)
		}
		return
	}
	discard := ctx.Resolver.Discard(ctx.Controller)
	candidates := discard
	if e.Type != "" {
		candidates = nil
		for _, id := range discard {
			if ctx.Resolver.TypeOf(id) == e.Type {
				candidates = append(candidates, id)
			}
		}
	}
	id, ok := ctx.ChooseCreature("Choose a "+e.noun()+" from your discard pile", candidates)
	if !ok {
		return
	}
	e.moveTo(ctx, id)
}

// DiscardHand discards cards from a player's hand: the chosen player's cards,
// optionally only creatures and only those of the house picked by an enclosing
// ChooseHouseThen. It models "discard each creature of the chosen house from your
// opponent's hand."
type DiscardHand struct {
	Player        Player
	CreaturesOnly bool
	OfChosenHouse bool
}

// validate rejects a DiscardHand whose player was left unset.
func (e DiscardHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardHand")
	}
	return nil
}

// Text renders the effect, e.g. "discard each creature of the chosen house from
// your opponent's hand".
func (e DiscardHand) Text() string {
	what := "each card"
	if e.CreaturesOnly {
		what = "each creature"
	}
	if e.OfChosenHouse {
		what += " of the chosen house"
	}
	whose := "your hand"
	if e.Player == Opponent {
		whose = "your opponent's hand"
	}
	return "discard " + what + " from " + whose
}

// Resolve discards every matching card from the chosen player's hand.
func (e DiscardHand) Resolve(ctx *EffectContext) {
	owner := ctx.PlayerFor(e.Player)
	for _, id := range ctx.Resolver.Hand(owner) {
		if e.CreaturesOnly && !ctx.Resolver.IsCreature(id) {
			continue
		}
		if e.OfChosenHouse && ctx.Resolver.House(id) != ctx.ChosenHouse {
			continue
		}
		ctx.Resolver.DiscardCardFromHand(owner, id)
	}
}
