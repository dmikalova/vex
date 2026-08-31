package engine

import "testing"

func TestEmitEnemyDestroyed(t *testing.T) {
	pileDef := NewCard("Pile", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterEnemyCreatureDestroyed, CaptureAember{
			Amount: 1,
			Target: Target{Kind: TargetChosenFriendlyCreature},
			Source: Opponent,
		}))

	t.Run("captures when an enemy creature is destroyed on your turn", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(pileDef, 0)
		friendly := g.AddToBattleline(testCreature("f", 5), 0)
		enemy := g.AddToBattleline(testCreature("e", 2), 1)
		g.State.Aember[1] = 3

		g.destroyEach(0, []LocalID{enemy})

		if g.AmberOn(friendly) != 1 {
			t.Errorf("friendly Æmber = %d, want 1", g.AmberOn(friendly))
		}
		if g.Aember(1) != 2 {
			t.Errorf("opponent pool = %d, want 2", g.Aember(1))
		}
	})

	t.Run("does not fire when a friendly creature is destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(pileDef, 0)
		friendly := g.AddToBattleline(testCreature("f", 5), 0)
		g.State.Aember[1] = 3

		g.destroyEach(0, []LocalID{friendly})

		if g.Aember(1) != 3 {
			t.Errorf("opponent pool = %d, want 3 (a friendly death does not trigger)", g.Aember(1))
		}
	})

	t.Run("does not fire on the opponent's turn", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 1
		g.AddArtifact(pileDef, 0)
		g.AddToBattleline(testCreature("f", 5), 0)
		mine := g.AddToBattleline(testCreature("m", 2), 0)
		g.State.Aember[1] = 3

		g.destroyEach(1, []LocalID{mine})

		if g.Aember(1) != 3 {
			t.Errorf("opponent pool = %d, want 3 (it is not your turn)", g.Aember(1))
		}
	})
}
