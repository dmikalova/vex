package engine

import "strings"

// PutFromHand puts a card the controller chooses from their own hand directly
// into play — Swap Widget swapping in a replacement creature. Type restricts the
// choice to cards of that type; House restricts it to that house; either's zero
// value allows any. ExceptSameName excludes a card sharing the name of the card
// currently in context (ctx.It), the "with a different name" clause — meant to
// follow a gate that left a card in context, e.g. Then{PutFromPlay, PutFromHand}.
type PutFromHand struct {
	Type           CardType
	House          House
	ExceptSameName bool
}

// noun renders the kind of card the effect puts into play, e.g. "Mars creature
// with a different name".
func (e PutFromHand) noun() string {
	base := "card"
	if e.Type != TypeUnset {
		base = strings.ToLower(e.Type.String())
	}
	if e.House != HouseNone {
		base = e.House.String() + " " + base
	}
	if e.ExceptSameName {
		base += " with a different name"
	}
	return base
}

// Text renders the effect, e.g. "put a Mars creature with a different name from
// your hand into play".
func (e PutFromHand) Text() string {
	return "put " + indefinite(e.noun()) + " from your hand into play"
}

// Resolve puts a card the controller chooses from their hand into play. Nothing
// happens if there is no candidate or the choice is declined. The chosen card is
// left in context (ctx.It) so a following effect can act on "it" (Swap Widget
// readying the creature it just put into play).
func (e PutFromHand) Resolve(ctx *EffectContext) {
	var candidates []LocalID
	for _, id := range ctx.Resolver.Hand(ctx.Controller) {
		if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
			continue
		}
		if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
			continue
		}
		if e.ExceptSameName && ctx.HasIt && ctx.Resolver.Name(id) == ctx.Resolver.Name(ctx.It) {
			continue
		}
		candidates = append(candidates, id)
	}
	id, ok := ctx.ChooseCard("Choose a "+e.noun()+" from your hand to put into play", candidates)
	if !ok {
		return
	}
	ctx.Resolver.PutIntoPlay(id, ctx.Controller)
	ctx.It, ctx.HasIt = id, true
}
