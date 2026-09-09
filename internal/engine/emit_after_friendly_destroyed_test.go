package engine

import "testing"

func TestEmitAfterFriendlyDestroyed(t *testing.T) {
	watcherDef := NewCard("Watcher", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterFriendlyCreatureDestroyed, GainAember{
			Player: Controller,
			Amount: 1,
		}))

	t.Run("fires when a creature under the same control is destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(watcherDef, 0)
		friendly := g.AddToBattleline(testCreature("f", 5), 0)

		g.destroyEach(0, []LocalID{friendly})

		if g.Aember(0) != 1 {
			t.Errorf("pool = %d, want 1", g.Aember(0))
		}
	})

	t.Run("does not fire when an enemy creature is destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(watcherDef, 0)
		enemy := g.AddToBattleline(testCreature("e", 2), 1)

		g.destroyEach(0, []LocalID{enemy})

		if g.Aember(0) != 0 {
			t.Errorf("pool = %d, want 0 (an enemy death does not trigger)", g.Aember(0))
		}
	})
}
