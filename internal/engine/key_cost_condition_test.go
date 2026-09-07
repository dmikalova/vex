package engine

import "testing"

// TestKeyCostChangeWhileCondition checks a key-cost change gated by a condition:
// Proclamation 346E taxes the opponent's keys only while they field creatures
// from fewer than three houses.
func TestKeyCostChangeWhileCondition(t *testing.T) {
	proclamation := func() CardDefinition {
		return NewCard("proc", Sanctum, Artifact, Rare,
			WithKeyCost(NewKeyCostChange(Opponent, 2).While(
				PlayerControlsFewerHousesThan{Player: Opponent, Amount: 3})))
	}

	t.Run("taxes while the opponent has fewer than three houses", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(proclamation(), 0)
		g.AddToBattleline(NewCard("a", Sanctum, Creature, Common, WithPower(1)), 1)
		g.AddToBattleline(NewCard("b", Untamed, Creature, Common, WithPower(1)), 1)
		if got := g.CurrentKeyCost(1); got != KeyCost+2 {
			t.Errorf("opponent key cost = %d, want %d", got, KeyCost+2)
		}
	})

	t.Run("stops taxing once the opponent fields three houses", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(proclamation(), 0)
		g.AddToBattleline(NewCard("a", Sanctum, Creature, Common, WithPower(1)), 1)
		g.AddToBattleline(NewCard("b", Untamed, Creature, Common, WithPower(1)), 1)
		g.AddToBattleline(NewCard("c", Logos, Creature, Common, WithPower(1)), 1)
		if got := g.CurrentKeyCost(1); got != KeyCost {
			t.Errorf("opponent key cost = %d, want %d", got, KeyCost)
		}
	})

	t.Run("never taxes its own controller", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(proclamation(), 0)
		if got := g.CurrentKeyCost(0); got != KeyCost {
			t.Errorf("own key cost = %d, want %d", got, KeyCost)
		}
	})

	t.Run("renders a while clause", func(t *testing.T) {
		want := "While your opponent does not control creatures from 3 or more " +
			"different houses, your opponent's keys cost +2 Æmber."
		if got := keyCostText(proclamation().KeyCostChanges[0]); got != want {
			t.Errorf("text = %q, want %q", got, want)
		}
	})
}

// TestPlayerControlsFewerHousesThanCondition exercises the condition's text and
// validation directly.
func TestPlayerControlsFewerHousesThanCondition(t *testing.T) {
	opp := PlayerControlsFewerHousesThan{Player: Opponent, Amount: 3}
	if got := opp.CondText(); got != "while your opponent does not control creatures from 3 or more different houses" {
		t.Errorf("CondText = %q", got)
	}
	you := PlayerControlsFewerHousesThan{Player: Controller, Amount: 2}
	if got := you.CondText(); got != "while you do not control creatures from 2 or more different houses" {
		t.Errorf("CondText = %q", got)
	}
	if err := (PlayerControlsFewerHousesThan{Player: Opponent}).validate(); err == nil {
		t.Error("zero amount should be rejected")
	}
	if err := opp.validate(); err != nil {
		t.Errorf("validate: %v", err)
	}
}
