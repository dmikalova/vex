package engine

import (
	"fmt"
	"strings"
)

// This file holds the effects that play a card out of a hand as part of resolving
// another card — the "play a card" clause, as opposed to a player taking their own
// play action.

// PlayFromHand has the controller play a card from their own hand right now,
// ignoring the active-house gate — Phase Shift's off-house card. House and Type
// narrow which cards may be chosen; Except inverts the house filter, so House names
// the house that may *not* be played ("a non-Logos card").
//
// KeyForge prints this as a permission held open for the rest of the turn ("you may
// play one non-Logos card this turn"); it is rendered and resolved as an immediate
// play instead, which needs no turn-scoped memory of an unspent allowance (see
// card-wording-rules.md rule 21). With no legal card in hand it does nothing.
//
//rulebook:effect Play a Card from Hand
type PlayFromHand struct {
	House  House
	Except bool
	Type   CardType
}

// validate rejects an inverted filter that names no house to exclude.
func (e PlayFromHand) validate() error {
	if e.Except && e.House == HouseNone {
		return fmt.Errorf("PlayFromHand: Except needs a house to exclude")
	}
	return nil
}

// Text renders the effect, e.g. "play a non-Logos card" or "play a Mars creature".
func (e PlayFromHand) Text() string { return "play " + indefinite(e.noun()) }

// noun names the cards the filters admit, e.g. "card", "non-Logos card", or
// "Mars creature".
func (e PlayFromHand) noun() string {
	noun := "card"
	if e.Type != TypeUnset && e.Type != AnyType {
		noun = strings.ToLower(e.Type.String())
	}
	switch {
	case e.Except:
		return "non-" + e.House.String() + " " + noun
	case e.House != HouseNone:
		return e.House.String() + " " + noun
	default:
		return noun
	}
}

// Resolve has the controller choose a matching card in hand and plays it.
func (e PlayFromHand) Resolve(ctx *EffectContext) {
	candidates := e.candidates(ctx)
	if len(candidates) == 0 {
		return
	}
	id, ok := ctx.ChooseCreature("Choose "+indefinite(e.noun())+" to play", candidates)
	if !ok {
		return
	}
	ctx.Resolver.PlayFromHand(ctx.Controller, id)
}

// candidates are the controller's hand cards the filters admit.
func (e PlayFromHand) candidates(ctx *EffectContext) []LocalID {
	var out []LocalID
	for _, id := range ctx.Resolver.Hand(ctx.Controller) {
		if e.Type != TypeUnset && e.Type != AnyType && ctx.Resolver.TypeOf(id) != e.Type {
			continue
		}
		if e.House != HouseNone && (ctx.Resolver.House(id) == e.House) == e.Except {
			continue
		}
		out = append(out, id)
	}
	return out
}
