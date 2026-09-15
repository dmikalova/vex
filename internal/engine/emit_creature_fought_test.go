package engine

import "testing"

func TestEmitCreatureFought(t *testing.T) {
	// Shattered-Throne-style: after any creature is used to fight, it captures 1
	// Æmber from its opponent.
	throneDef := NewCard(
		"Throne",
		Brobnar,
		Artifact,
		Uncommon,
		WithAbility(
			TriggerAfterCreatureFights,
			CaptureAember{
				Amount: 1,
				Target: Target{Kind: TargetTriggeringCreature},
				Source: ItsOpponent,
			},
		),
	)

	t.Run("the fighting creature captures after fighting", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.State.Aember[1] = 3
		g.AddArtifact(throneDef, 0)
		attacker := g.AddToBattleline(testCreature("att", 5), 0)
		defender := g.AddToBattleline(testCreature("def", 3), 1)

		if err := g.Fight(0, attacker, defender); err != nil {
			t.Fatalf("fight: %v", err)
		}
		if got := g.AmberOn(attacker); got != 1 {
			t.Errorf("captured Æmber on attacker = %d, want 1", got)
		}
		if g.State.Aember[1] != 2 {
			t.Errorf("opponent pool = %d, want 2", g.State.Aember[1])
		}
	})

	t.Run("fires for both players' throne", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.State.Aember[1] = 3
		g.AddArtifact(throneDef, 1) // enemy-controlled throne still fires
		attacker := g.AddToBattleline(testCreature("att", 5), 0)
		defender := g.AddToBattleline(testCreature("def", 3), 1)

		if err := g.Fight(0, attacker, defender); err != nil {
			t.Fatalf("fight: %v", err)
		}
		if got := g.AmberOn(attacker); got != 1 {
			t.Errorf("captured Æmber on attacker = %d, want 1", got)
		}
	})
}

// TestAfterFriendlyCreatureFights covers the friendly-scoped fight trigger
// (Lieutenant Gorvenal): it fires when a creature on the source's own side fights,
// and not when an enemy creature fights.
func TestAfterFriendlyCreatureFights(t *testing.T) {
	watcherDef := NewCard("Watcher", Sanctum, Creature, Common, WithPower(4),
		WithAbility(TriggerAfterFriendlyCreatureFights,
			CaptureAember{Amount: 1, Target: Target{Kind: TargetThisCreature}, Source: Opponent}))

	t.Run("fires when a friendly creature fights", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.State.Aember[1] = 3
		watcher := g.AddToBattleline(watcherDef, 0)
		ally := g.AddToBattleline(testCreature("ally", 5), 0)
		defender := g.AddToBattleline(testCreature("def", 3), 1)

		if err := g.Fight(0, ally, defender); err != nil {
			t.Fatalf("fight: %v", err)
		}
		if got := g.AmberOn(watcher); got != 1 {
			t.Errorf("captured on watcher = %d, want 1 after a friendly fight", got)
		}
	})

	t.Run("does not fire when an enemy creature fights", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 1
		g.State.Aember[0] = 3
		watcher := g.AddToBattleline(watcherDef, 0)
		enemy := g.AddToBattleline(testCreature("enemy", 5), 1)
		defender := g.AddToBattleline(testCreature("def", 3), 0)

		if err := g.Fight(1, enemy, defender); err != nil {
			t.Fatalf("fight: %v", err)
		}
		if got := g.AmberOn(watcher); got != 0 {
			t.Errorf("watcher captured %d after an enemy fight, want 0", got)
		}
	})
}
