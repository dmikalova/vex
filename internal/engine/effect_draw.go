package engine

import "fmt"

// Drawing puts the top card of your deck into your hand. If your deck is empty
// when you must draw, your discard pile is shuffled to form a new deck first, so
// you only fail to draw when both deck and discard are empty.
type Draw struct {
	// Amount is how many cards to draw; Per multiplies it by a running count.
	Amount int
	Per    Count
	// Or switches Amount to an alternate when a condition holds, so the card reads
	// "draw a card, or 2 cards if …" instead of a two-armed Otherwise branch (Hyde
	// draws 2 while controlling Velum).
	Or OrAmount
	// You names the drawer explicitly ("you draw a card") to re-assert the subject
	// when the draw follows a clause whose subject was the opponent (Perplexing
	// Sophistry: "your opponent discards ..., and you draw a card").
	You bool
}

// validate checks the Or guard, if one is set.
func (e Draw) validate() error { return e.Or.validate() }

// drawObject renders the noun the draw acts on, e.g. "a card" or "2 cards".
func drawObject(n int) string {
	if n == 1 {
		return "a card"
	}
	return fmt.Sprintf("%d cards", n)
}

// Text renders the effect, e.g. "draw a card" or "draw 2 cards". A "for each"
// count leads the sentence (rule 9), e.g. "for each Mars card in your hand, draw a
// card".
func (e Draw) Text() string {
	phrase := "draw " + drawObject(e.Amount)
	if e.You {
		phrase = "you " + phrase
	}
	body := forEach(e.Per, phrase)
	if e.Or.set() {
		body += e.Or.tail(drawObject(e.Or.Amount))
	}
	return body
}

// Resolve draws the cards, scaling by the Per count when one is set.
func (e Draw) Resolve(ctx *EffectContext) {
	ctx.Resolver.Draw(ctx.Controller, scaled(e.Or.pick(e.Amount, ctx), e.Per, ctx), ctx.Source)
}
