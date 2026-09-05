package engine

import "fmt"

// Drawing puts the top card of your deck into your hand. If your deck is empty
// when you must draw, your discard pile is shuffled to form a new deck first, so
// you only fail to draw when both deck and discard are empty.
type Draw struct {
	// Amount is how many cards to draw; Per multiplies it by a running count.
	Amount int
	Per    Count
	// You names the drawer explicitly ("you draw a card") to re-assert the subject
	// when the draw follows a clause whose subject was the opponent (Perplexing
	// Sophistry: "your opponent discards ..., and you draw a card").
	You bool
}

// Text renders the effect, e.g. "draw a card" or "draw 2 cards". A "for each"
// count leads the sentence (rule 9), e.g. "for each Mars card in your hand, draw a
// card".
func (e Draw) Text() string {
	phrase := fmt.Sprintf("draw %d cards", e.Amount)
	if e.Amount == 1 {
		phrase = "draw a card"
	}
	if e.You {
		phrase = "you " + phrase
	}
	return forEach(e.Per, phrase)
}

// Resolve draws the cards, scaling by the Per count when one is set.
func (e Draw) Resolve(ctx *EffectContext) {
	ctx.Resolver.Draw(ctx.Controller, scaled(e.Amount, e.Per, ctx))
}
