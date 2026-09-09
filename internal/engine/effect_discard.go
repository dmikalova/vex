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
	// Name restricts the choice to cards with that exact name; the zero value allows
	// any (Ortannu the Chained returns each copy of Ortannu's Binding).
	Name string
	// Destination is where the card goes: ToHand or ToTopOfDeck.
	Destination Destination
	// All moves every matching card instead of one chosen card (Arise! returning
	// each creature of a house).
	All bool
	// OfChosenHouse limits the matching cards to the house an enclosing
	// ChooseHouseThen picked. It applies only with All.
	OfChosenHouse bool
}

// noun renders the kind of card the effect moves — the card's own name when Name
// is set (e.g. "Ortannu's Binding"), the lowercased card type when Type is set
// (e.g. "creature"), otherwise the generic "card".
func (e PutFromDiscard) noun() string {
	if e.Name != "" {
		return e.Name
	}
	base := "card"
	if e.Type != TypeUnset {
		base = strings.ToLower(e.Type.String())
	}
	if e.Trait != traitUnset {
		base = e.Trait.String() + " " + base
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

// admits reports whether a discard-pile card passes the Type / Trait / Name
// filters (the OfChosenHouse filter is applied separately, only with All).
func (e PutFromDiscard) admits(ctx *EffectContext, id LocalID) bool {
	if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
		return false
	}
	if e.Trait != traitUnset && !ctx.Resolver.HasTrait(id, e.Trait) {
		return false
	}
	return e.Name == "" || ctx.Resolver.Name(id) == e.Name
}

// Resolve moves a card from the controller's discard pile to the destination. With
// All it moves every matching card; otherwise the controller chooses one, and
// nothing happens if there is no candidate or the choice is declined.
func (e PutFromDiscard) Resolve(ctx *EffectContext) {
	if e.All {
		for _, id := range discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
			return e.admits(ctx, id) &&
				(!e.OfChosenHouse || ctx.Resolver.House(id) == ctx.ChosenHouse)
		}) {
			e.moveTo(ctx, id)
		}
		return
	}
	candidates := discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		return e.admits(ctx, id)
	})
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
	// Player names whose hand is discarded from.
	Player Player
	// Types restricts the discard to cards of the listed types; empty discards any card.
	Types []CardType
	// OfChosenHouse limits the discard to cards of the house picked by an enclosing
	// ChooseHouseThen.
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
	for _, id := range handCardsWhere(ctx, owner, func(id LocalID) bool {
		if !matchesTypes(e.Types, ctx.Resolver.TypeOf(id)) {
			return false
		}
		return !e.OfChosenHouse || ctx.Resolver.House(id) == ctx.ChosenHouse
	}) {
		ctx.Resolver.DiscardCardFromHand(owner, id)
	}
}

// DiscardRandomFromHand discards one or more uniformly random cards from a player's
// hand — the "discard a random card" effect on cards like Mind Barb and Tocsin,
// where the discarding player does not choose which card leaves a hidden hand.
// Amount discards that many (Nogi Smartfist discards 2); an unset zero means one.
type DiscardRandomFromHand struct {
	Player Player
	Amount int
}

// count is how many cards to discard: the named Amount, or one when left unset.
func (e DiscardRandomFromHand) count() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// validate rejects a DiscardRandomFromHand whose player was left unset.
func (e DiscardRandomFromHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardRandomFromHand")
	}
	return nil
}

// Text renders the effect, e.g. "your opponent discards a random card from their
// hand" or "discard 2 random cards from your hand".
func (e DiscardRandomFromHand) Text() string {
	object := "a random card"
	if e.count() > 1 {
		object = fmt.Sprintf("%d random cards", e.count())
	}
	switch e.Player {
	case Opponent:
		return "your opponent discards " + object + " from their hand"
	case ItsOwner:
		return "its owner discards " + object + " from their hand"
	default:
		return "discard " + object + " from your hand"
	}
}

// Resolve discards count() random cards from the chosen player's hand, one at a
// time (a smaller hand simply discards fewer).
func (e DiscardRandomFromHand) Resolve(ctx *EffectContext) {
	p := ctx.PlayerFor(e.Player)
	for i := 0; i < e.count(); i++ {
		ctx.Resolver.DiscardRandomFromHand(p)
	}
}

// DiscardRandomFromArchives discards one uniformly random card from a player's
// archives — Tantadlin's "discard a random card from your opponent's archives",
// where the discarding player cannot see the facedown archives to choose.
type DiscardRandomFromArchives struct {
	Player Player
}

// validate rejects a DiscardRandomFromArchives whose player was left unset.
func (e DiscardRandomFromArchives) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardRandomFromArchives")
	}
	return nil
}

// Text renders the effect, e.g. "discard a random card from your opponent's
// archives".
func (e DiscardRandomFromArchives) Text() string {
	switch e.Player {
	case Opponent:
		return "discard a random card from your opponent's archives"
	case ItsOwner:
		return "its owner discards a random card from their archives"
	default:
		return "discard a random card from your archives"
	}
}

// Resolve discards one random card from the chosen player's archives.
func (e DiscardRandomFromArchives) Resolve(ctx *EffectContext) {
	ctx.Resolver.DiscardRandomFromArchives(ctx.PlayerFor(e.Player))
}

// DiscardFromHand has the controller choose and discard Amount cards from their own
// hand — the "discard a card" effect where the player picks which card leaves
// (Sloppy Labwork), distinct from DiscardHand (which discards every matching card)
// and DiscardRandomFromHand (which the player does not choose). Types limits the
// choice to the listed card types (Feeding Pit's "discard a creature from your
// hand"); an empty Types allows any card. AnyNumber lets the controller discard as
// many matching cards as they like, declining when done — Helmsman Spears discards
// any number of cards. Every card discarded this way is recorded on the context so
// a following ForEachDiscarded can act once per card.
type DiscardFromHand struct {
	Amount    int
	Types     []CardType
	AnyNumber bool
}

// Text renders the effect, e.g. "discard a card from your hand", naming the source
// zone explicitly (rule 17).
func (e DiscardFromHand) Text() string {
	noun := typeNoun(e.Types)
	if e.AnyNumber {
		return "discard any number of " + noun + "s from your hand"
	}
	if e.Amount == 1 {
		return "discard " + indefinite(noun) + " from your hand"
	}
	return "discard " + countNoun(e.Amount, noun) + " from your hand"
}

// Resolve has the controller choose and discard Amount cards from their hand,
// stopping early if the hand runs out or the choice is declined.
func (e DiscardFromHand) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate performs the discards and reports whether any card was discarded, so
// DiscardFromHand can gate a Then — Feeding Pit only gains Æmber if a creature was
// discarded. Each discarded card is appended to ctx.Produced.Discarded.
func (e DiscardFromHand) resolveGate(ctx *EffectContext) bool {
	moved := false
	for i := 0; e.AnyNumber || i < e.Amount; i++ {
		candidates := handCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
			return matchesTypes(e.Types, ctx.Resolver.TypeOf(id))
		})
		if len(candidates) == 0 {
			return moved
		}
		var (
			id LocalID
			ok bool
		)
		if e.AnyNumber {
			id, ok = ctx.ChooseCardOptional("Choose a card to discard", candidates)
		} else {
			id, ok = ctx.ChooseCreature("Choose a card to discard", candidates)
		}
		if !ok {
			return moved
		}
		ctx.Resolver.DiscardCardFromHand(ctx.Controller, id)
		ctx.Produced.Discarded = append(ctx.Produced.Discarded, id)
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
