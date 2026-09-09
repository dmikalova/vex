package engine

import "testing"

// TestEmitAfterEnemyDestroyedFighting covers The Colosseum's trigger: a bystander
// artifact reacts when an enemy creature is destroyed in a fight, whether the
// enemy dies as the defender or its Hazardous kills the attacker, and stays quiet
// when the creature destroyed is friendly to it.
func TestEmitAfterEnemyDestroyedFighting(t *testing.T) {
	watcher := func() CardDefinition {
		return NewCard(
			"Colosseum",
			Saurian,
			Artifact,
			Rare,
			WithAbility(
				TriggerAfterEnemyDestroyedFighting,
				GainAember{Player: Controller, Amount: 1},
			),
		)
	}

	t.Run("fires when the enemy defender is destroyed in the fight", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(watcher(), 0)
		hunter := g.AddToBattleline(testCreature("hunter", 6), 0)
		prey := g.AddToBattleline(testCreature("prey", 1), 1)
		before := g.Aember(0)
		if err := g.Fight(0, hunter, prey); err != nil {
			t.Fatalf("Fight: %v", err)
		}
		if g.Aember(0) != before+1 {
			t.Errorf("bystander controller aember = %d, want %d", g.Aember(0), before+1)
		}
	})

	t.Run("fires for the defender's side when Hazardous kills the attacker", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(watcher(), 1) // watcher controlled by the defender's side
		weak := g.AddToBattleline(testCreature("weak", 1), 0)
		guard := g.AddToBattleline(NewCard("guard", Sanctum, Creature, Common,
			WithPower(6), WithHazardous(5)), 1)
		before := g.Aember(1)
		if err := g.Fight(0, weak, guard); err != nil {
			t.Fatalf("Fight: %v", err)
		}
		if g.inPlay(weak) {
			t.Fatal("the attacker should be destroyed by Hazardous")
		}
		if g.Aember(1) != before+1 {
			t.Errorf("defender-side bystander aember = %d, want %d", g.Aember(1), before+1)
		}
	})

	t.Run("stays quiet when a friendly creature dies in the fight", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(watcher(), 0) // watches for enemy deaths only
		weak := g.AddToBattleline(testCreature("weak", 1), 0)
		guard := g.AddToBattleline(NewCard("guard", Sanctum, Creature, Common,
			WithPower(6), WithHazardous(5)), 1)
		before := g.Aember(0)
		if err := g.Fight(0, weak, guard); err != nil {
			t.Fatalf("Fight: %v", err)
		}
		if g.Aember(0) != before {
			t.Errorf(
				"bystander aember = %d, want %d (a friendly death does not fire)",
				g.Aember(0),
				before,
			)
		}
	})
}
