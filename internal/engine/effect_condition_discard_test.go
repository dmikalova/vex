package engine

import "testing"

// TestCardsInDiscardAtLeast covers the threshold Low Dawn reads: 3 or more
// Untamed creatures in the controller's discard pile.
func TestCardsInDiscardAtLeast(t *testing.T) {
	if got := (CardsInDiscardAtLeast{House: namedHouse(Untamed), Type: Creature, Amount: 3}).
		CondText(); got != "if there are 3 or more Untamed creatures in your discard pile" {
		t.Errorf("text = %q", got)
	}
	if (CardsInDiscardAtLeast{House: namedHouse(Untamed), Type: Creature}).validate() == nil {
		t.Error("a zero Amount should be invalid")
	}
	if err := (CardsInDiscardAtLeast{Amount: 1}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	e := CardsInDiscardAtLeast{House: namedHouse(Untamed), Type: Creature, Amount: 3}

	// Two Untamed creatures fall short of the threshold, even alongside a
	// non-matching Untamed tactic and a Mars creature.
	g := NewGame("A", "B", 1)
	g.AddToDiscard(NewCard("u1", Untamed, Creature, Common, WithPower(3)), 0)
	g.AddToDiscard(NewCard("u2", Untamed, Creature, Common, WithPower(3)), 0)
	g.AddToDiscard(NewCard("ut", Untamed, Tactic, Common), 0)
	g.AddToDiscard(NewCard("mc", Mars, Creature, Common, WithPower(3)), 0)
	if (CardsInDiscardAtLeast{House: namedHouse(Untamed), Type: Creature, Amount: 3}).
		Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("two Untamed creatures should not meet the threshold of 3")
	}

	// A third Untamed creature reaches exactly the threshold.
	g.AddToDiscard(NewCard("u3", Untamed, Creature, Common, WithPower(3)), 0)
	if !e.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("three Untamed creatures should meet the threshold of 3")
	}
}

// TestNamedCardInDiscard covers the Monuments' gate: whether a card of a given
// name waits in the controller's discard pile.
func TestNamedCardInDiscard(t *testing.T) {
	e := NamedCardInDiscard{Name: "Faust the Great"}
	if got := e.CondText(); got != "if Faust the Great is in your discard pile" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	// A card of the name in the opponent's discard pile does not satisfy the
	// controller's condition, nor does a differently named card in their own.
	g.AddToDiscard(NewCard("Faust the Great", Saurian, Creature, Common, WithPower(4)), 1)
	g.AddToDiscard(NewCard("Cornicen Octavia", Saurian, Creature, Common, WithPower(5)), 0)
	if e.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("Faust in the opponent's discard should not meet the condition")
	}

	// The named card in the controller's own discard pile satisfies it.
	g.AddToDiscard(NewCard("Faust the Great", Saurian, Creature, Common, WithPower(4)), 0)
	if !e.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("Faust in your own discard should meet the condition")
	}
}
