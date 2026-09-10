package engine

import "testing"

// TestItIsNotOfHouse covers the "non-<house> card" condition: it is met only when
// a card is in context and its house differs from the named house.
func TestItIsNotOfHouse(t *testing.T) {
	c := ItIsNotOfHouse{House: StarAlliance}
	if got := c.CondText(); got != "if it is a non-Star Alliance card" {
		t.Errorf("CondText = %q", got)
	}

	g := started(t)
	star := g.AddToDeck(NewCard("marine", StarAlliance, Creature, Common, WithPower(2)), 0)
	other := g.AddToDeck(NewCard("giant", Brobnar, Creature, Common, WithPower(2)), 0)

	// A card of another house is a non-Star Alliance card.
	if !c.Met(&EffectContext{Resolver: g, Controller: 0, It: other, HasIt: true}) {
		t.Error("a Brobnar card should be a non-Star Alliance card")
	}
	// A Star Alliance card is not.
	if c.Met(&EffectContext{Resolver: g, Controller: 0, It: star, HasIt: true}) {
		t.Error("a Star Alliance card should not be a non-Star Alliance card")
	}
	// With no card in context the condition is not met.
	if c.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("no card in context should not meet the condition")
	}
}
