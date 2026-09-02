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
	// Trait restricts the choice to cards with that trait; the zero value allows any.
	Trait Trait
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
	base := "card"
	if e.Type != TypeUnset {
		base = strings.ToLower(e.Type.String())
	}
	if e.Trait != "" {
		base = string(e.Trait) + " trait " + base
	}
	return base
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
			if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
				continue
			}
			if e.Trait != "" && !ctx.Resolver.HasTrait(id, e.Trait) {
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
	if e.Type != TypeUnset || e.Trait != "" {
		candidates = nil
		for _, id := range discard {
			if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
				continue
			}
			if e.Trait != "" && !ctx.Resolver.HasTrait(id, e.Trait) {
				continue
			}
			candidates = append(candidates, id)
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
	Player Player
	// Types restricts the discard to cards of the listed types; empty discards any card.
	Types         []CardType
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
	what := "each " + typeNoun(e.Types)
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
		if !matchesTypes(e.Types, ctx.Resolver.TypeOf(id)) {
			continue
		}
		if e.OfChosenHouse && ctx.Resolver.House(id) != ctx.ChosenHouse {
			continue
		}
		ctx.Resolver.DiscardCardFromHand(owner, id)
	}
}

// DiscardRandomFromHand discards one uniformly random card from a player's hand — the
// "discard a random card" effect on cards like Mind Barb and Tocsin, where the
// discarding player does not choose which card leaves a hidden hand.
type DiscardRandomFromHand struct {
	Player Player
}

// validate rejects a DiscardRandomFromHand whose player was left unset.
func (e DiscardRandomFromHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardRandomFromHand")
	}
	return nil
}

// Text renders the effect, e.g. "your opponent discards a random card from their
// hand".
func (e DiscardRandomFromHand) Text() string {
	switch e.Player {
	case Opponent:
		return "your opponent discards a random card from their hand"
	case ItsOwner:
		return "its owner discards a random card from their hand"
	default:
		return "discard a random card from your hand"
	}
}

// Resolve discards one random card from the chosen player's hand.
func (e DiscardRandomFromHand) Resolve(ctx *EffectContext) {
	ctx.Resolver.DiscardRandomFromHand(ctx.PlayerFor(e.Player))
}

// DiscardFromHand has the controller choose and discard Count cards from their own
// hand — the "discard a card" effect where the player picks which card leaves
// (Sloppy Labwork), distinct from DiscardHand (which discards every matching card)
// and DiscardRandomFromHand (which the player does not choose). Types limits the
// choice to the listed card types (Feeding Pit's "discard a creature from your
// hand"); an empty Types allows any card.
type DiscardFromHand struct {
	Count int
	Types []CardType
}

// Text renders the effect, e.g. "discard a card from your hand", naming the source
// zone explicitly (rule 17).
func (e DiscardFromHand) Text() string {
	noun := typeNoun(e.Types)
	if e.Count == 1 {
		return "discard a " + noun + " from your hand"
	}
	return fmt.Sprintf("discard %d %ss from your hand", e.Count, noun)
}

// Resolve has the controller choose and discard Count cards from their hand,
// stopping early if the hand runs out or the choice is declined.
func (e DiscardFromHand) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate performs the discards and reports whether any card was discarded, so
// DiscardFromHand can gate a Then — Feeding Pit only gains Æmber if a creature was
// discarded.
func (e DiscardFromHand) resolveGate(ctx *EffectContext) bool {
	moved := false
	for i := 0; i < e.Count; i++ {
		var candidates []LocalID
		for _, id := range ctx.Resolver.Hand(ctx.Controller) {
			if matchesTypes(e.Types, ctx.Resolver.TypeOf(id)) {
				candidates = append(candidates, id)
			}
		}
		if len(candidates) == 0 {
			return moved
		}
		id, ok := ctx.ChooseCreature("Choose a card to discard", candidates)
		if !ok {
			return moved
		}
		ctx.Resolver.DiscardCardFromHand(ctx.Controller, id)
		moved = true
	}
	return moved
}

// matchesTypes reports whether a card of type t passes a type filter: an empty
// filter allows any card, otherwise t must be one of the listed types.
func matchesTypes(types []CardType, t CardType) bool {
	if len(types) == 0 {
		return true
	}
	for _, want := range types {
		if want == t {
			return true
		}
	}
	return false
}

// typeNoun renders a card-type filter as a noun for discard text — the type nouns
// joined with "or" ("creature", "creature or artifact"), or "card" for no filter.
func typeNoun(types []CardType) string {
	if len(types) == 0 {
		return "card"
	}
	nouns := make([]string, len(types))
	for i, t := range types {
		nouns[i] = cardTypeNoun(t)
	}
	return strings.Join(nouns, " or ")
}

// cardTypeNoun renders a single card type as its discard noun.
func cardTypeNoun(t CardType) string {
	switch t {
	case Artifact:
		return "artifact"
	case Creature:
		return "creature"
	default:
		return "card"
	}
}
