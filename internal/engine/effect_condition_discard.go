package engine

import "fmt"

// CardsInDiscardAtLeast is met when the controller's discard pile holds at least
// Amount cards its Filter admits — Low Dawn gains Æmber only while 3 or more
// Untamed creatures wait in the discard pile. The zero Filter counts every card.
type CardsInDiscardAtLeast struct {
	// Filter narrows which cards in the pile count and supplies the noun the clause
	// prints. A card in a pile is in no battleline, so validate rejects an in-play
	// axis rather than letting it silently match nothing.
	Filter Filter
	Amount int
}

// validate requires a positive threshold and a filter a card in a pile can answer.
func (e CardsInDiscardAtLeast) validate() error {
	if e.Amount <= 0 {
		return fmt.Errorf("CardsInDiscardAtLeast: Amount must be positive")
	}
	return e.Filter.validateIdentityOnly("CardsInDiscardAtLeast")
}

// CondText renders the condition, e.g. "if there are 3 or more Untamed creatures
// in your discard pile".
func (e CardsInDiscardAtLeast) CondText() string {
	return fmt.Sprintf("if there are %d or more %s in your discard pile",
		e.Amount, plural(e.Amount, e.Filter.noun("card")))
}

// Met reports whether at least Amount matching cards sit in the controller's
// discard pile.
func (e CardsInDiscardAtLeast) Met(ctx *EffectContext) bool {
	matched := discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		return e.Filter.matches(ctx, id)
	})
	return len(matched) >= e.Amount
}

// DiscardedThisWay is met when the current dig or discard recorded at least one
// card its Filter admits — Saurian Egg checks whether it discarded any Saurian
// creature this way before it hatches and destroys itself. The zero Filter admits
// every discarded card.
type DiscardedThisWay struct {
	// Filter narrows which discarded cards count and supplies the noun the clause
	// prints. A discarded card is in no battleline, so validate rejects an in-play
	// axis.
	Filter Filter
}

// validate rejects a filter that reads the board, which a discarded card is not on.
func (e DiscardedThisWay) validate() error {
	return e.Filter.validateIdentityOnly("DiscardedThisWay")
}

// CondText renders the condition, e.g. "if you discard a Saurian creature this
// way".
func (e DiscardedThisWay) CondText() string {
	return "if you discard " + e.Filter.object("card") + " this way"
}

// Met reports whether any card the current dig or discard recorded is one the
// filter admits.
func (e DiscardedThisWay) Met(ctx *EffectContext) bool {
	for _, id := range discardedThisWay(ctx) {
		if e.Filter.matches(ctx, id) {
			return true
		}
	}
	return false
}

// NamedCardInDiscard is met when a card of a given name sits in the controller's
// discard pile — the Monuments strengthen their action when their namesake
// creature (Faust the Great, Cornicen Octavia, Consul Primus, Citizen Shrix) waits
// in the discard pile. It names the card by its printed name, not the source.
type NamedCardInDiscard struct {
	// Filter names the card to look for. It is a Filter rather than a bare name so
	// a card asking for "a Mars creature in your discard pile" has a field rather
	// than a new node; a filter naming a card prints it outright, without an
	// article. A card in a pile is in no battleline, so validate rejects an in-play
	// axis.
	Filter Filter
}

// validate rejects a filter that reads the board, which a pile card is not on.
func (c NamedCardInDiscard) validate() error {
	return c.Filter.validateIdentityOnly("NamedCardInDiscard")
}

// CondText renders the condition naming the card it looks for.
func (c NamedCardInDiscard) CondText() string {
	return "if " + c.Filter.object("card") + " is in your discard pile"
}

// Met reports whether a card the filter admits sits in the controller's discard
// pile.
func (c NamedCardInDiscard) Met(ctx *EffectContext) bool {
	for _, id := range ctx.Resolver.Discard(ctx.Controller) {
		if c.Filter.matches(ctx, id) {
			return true
		}
	}
	return false
}
