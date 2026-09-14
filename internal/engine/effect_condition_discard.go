package engine

import "fmt"

// CardsInDiscardAtLeast is met when the controller's discard pile holds at least
// Amount cards matching House and Type — Low Dawn gains Æmber only while 3 or more
// Untamed creatures wait in the discard pile. An unset House or Type applies no
// filter on that axis.
type CardsInDiscardAtLeast struct {
	House  HouseMatcher
	Type   CardType
	Amount int
}

// validate requires a positive threshold.
func (e CardsInDiscardAtLeast) validate() error {
	if e.Amount <= 0 {
		return fmt.Errorf("CardsInDiscardAtLeast: Amount must be positive")
	}
	return nil
}

// CondText renders the condition, e.g. "if there are 3 or more Untamed creatures
// in your discard pile".
func (e CardsInDiscardAtLeast) CondText() string {
	return fmt.Sprintf("if there are %d or more %s in your discard pile",
		e.Amount, plural(e.Amount, e.House.qualifyNoun(typeNoun(e.Type))))
}

// Met reports whether at least Amount matching cards sit in the controller's
// discard pile.
func (e CardsInDiscardAtLeast) Met(ctx *EffectContext) bool {
	matched := discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		if !e.House.matches(ctx, id) {
			return false
		}
		return e.Type == TypeUnset || ctx.Resolver.TypeOf(id) == e.Type
	})
	return len(matched) >= e.Amount
}
