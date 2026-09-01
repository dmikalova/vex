package engine

import (
	"testing"
)

func TestCardsDiscarded(t *testing.T) {
	e := CardsDiscarded{Player: Controller, House: Untamed, Amount: 1}
	if got := e.CondText(); got != "if you have discarded an Untamed card from your hand this turn" {
		t.Errorf("CondText() = %q", got)
	}
	eOpp := CardsDiscarded{Player: Opponent, House: Mars, Amount: 1}
	if got := eOpp.CondText(); got != "if your opponent has discarded a Mars card from their hand this turn" {
		t.Errorf("CondText() = %q", got)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if e.Met(ctx) {
		t.Error("Met() should be false initially")
	}

	if eOpp.Met(ctx) {
		t.Error("eOpp.Met() should be false initially")
	}

	if e.validate() != nil {
		t.Errorf("validate() = %v, want nil", e.validate())
	}
	if (CardsDiscarded{Player: Controller, House: Untamed}).validate() == nil {
		t.Error("validate() should reject a zero Amount")
	}
	eAmount := CardsDiscarded{Player: Controller, House: Untamed, Amount: 3}
	if got := eAmount.CondText(); got != "if you have discarded 3 Untamed cards from your hand this turn" {
		t.Errorf("eAmount.CondText() = %q", got)
	}

	c := g.AddToHand(NewCard("Untamed Action", Untamed, Tactic, Common), 0)
	g.DiscardCardFromHand(0, c)

	if !e.Met(ctx) {
		t.Error("Met() should be true after discarding Untamed card")
	}

	if got := g.DiscardedThisTurn(0); len(got) != 1 || g.House(got[0]) != Untamed {
		t.Errorf("DiscardedThisTurn = %v, want one Untamed card", got)
	}
}

func TestNewCardRejectsZeroAmountUseCondition(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject a UseCondition with a non-positive Amount")
		}
	}()
	NewCard("bad", Untamed, Creature, Rare, WithRestrictions(Restrictions{
		UseCondition: CardsDiscarded{Player: Controller, House: Untamed},
	}))
}
