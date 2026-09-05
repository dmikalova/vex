package engine

import "fmt"

// Drawing puts the top card of your deck into your hand. If your deck is empty
// when you must draw, your discard pile is shuffled to form a new deck first, so
// you only fail to draw when both deck and discard are empty.
type Draw struct {
	Amount int
	Per    Count
}

// Text renders the effect, e.g. "draw a card" or "draw 2 cards". A "for each"
// count leads the sentence (rule 9), e.g. "for each Mars card in your hand, draw a
// card".
func (e Draw) Text() string {
	phrase := fmt.Sprintf("draw %d cards", e.Amount)
	if e.Amount == 1 {
		phrase = "draw a card"
	}
	return forEach(e.Per, phrase)
}

// Resolve draws the cards, scaling by the Per count when one is set.
func (e Draw) Resolve(ctx *EffectContext) {
	amount := e.Amount
	if e.Per != nil {
		amount *= e.Per.Value(ctx)
	}
	ctx.Resolver.Draw(ctx.Controller, amount)
}
