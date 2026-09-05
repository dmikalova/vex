package engine

import "testing"

func TestEmitCreatureReaped(t *testing.T) {
	// Orb-style: after any creature reaps, stun it.
	orbDef := NewCard(
		"Orb",
		Dis,
		Artifact,
		Rare,
		WithAbility(
			TriggerAfterCreatureReaps,
			Stun{Target: Target{Kind: TargetTriggeringCreature}},
		),
	)
	// Pip-style: after an enemy creature reaps, stun it.
	pipDef := NewCard(
		"Pip",
		Logos,
		Creature,
		Common,
		WithAbility(
			TriggerAfterEnemyCreatureReaps,
			Stun{Target: Target{Kind: TargetTriggeringCreature}},
		),
	)

	t.Run("enemy Orb stuns the reaper", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(orbDef, 1)
		reaper := g.AddToBattleline(testCreature("r", 3), 0)

		g.reapWith(reaper)

		if !g.State.Cards[reaper].Stunned {
			t.Error("reaper should be stunned by the enemy Orb")
		}
	})

	t.Run("friendly Orb stuns your own reaper too", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddArtifact(orbDef, 0)
		reaper := g.AddToBattleline(testCreature("r", 3), 0)

		g.reapWith(reaper)

		if !g.State.Cards[reaper].Stunned {
			t.Error("friendly Orb should stun your own reaper (any creature)")
		}
	})

	t.Run("enemy-reap reaction fires for the reaper's opponent", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddToBattleline(pipDef, 1)
		reaper := g.AddToBattleline(testCreature("r", 3), 0)

		g.reapWith(reaper)

		if !g.State.Cards[reaper].Stunned {
			t.Error("the enemy reaper should be stunned by Pip")
		}
	})

	t.Run("enemy-reap reaction does not fire on the controller's own reap", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.State.ActivePlayer = 0
		g.AddToBattleline(pipDef, 0)
		reaper := g.AddToBattleline(testCreature("r", 3), 0)

		g.reapWith(reaper)

		if g.State.Cards[reaper].Stunned {
			t.Error("your own reap should not fire your Pip's enemy-reap reaction")
		}
	})
}
