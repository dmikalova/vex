package engine

import "testing"

// A constant ability with DisableTriggers stops the listed triggers from firing
// board-wide while its source stays in play: Purifier of Souls disables every
// Destroyed ability, so a creature destroyed alongside it gains its controller
// nothing.
func TestDisableTriggersStopsDestroyedAbilities(t *testing.T) {
	g := started(t)
	// A buffer with a plain constant ability (no DisableTriggers) is in play too,
	// so the disable scan skips a source that disables nothing.
	g.AddArtifact(NewCard("buffer", Dis, Artifact, Rare,
		WithConstantAbility(ConstantAbility{PowerBonus: 1})), 0)
	// The purifier lists a non-matching trigger before Destroyed, so the scan walks
	// past a trigger it does not disable before finding the one it does.
	g.AddArtifact(NewCard("purifier", Sanctum, Artifact, Rare,
		WithConstantAbility(ConstantAbility{
			DisableTriggers: []Trigger{TriggerAction, TriggerDestroyed},
		})), 0)
	victim := g.AddToBattleline(NewCard("v", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerDestroyed, GainAember{Player: Controller, Amount: 1})), 0)

	g.DestroyEach(0, []LocalID{victim})

	if g.Aember(0) != 0 {
		t.Errorf("aember = %d, want 0 (Destroyed ability must not fire)", g.Aember(0))
	}
}

func TestDisableTriggersText(t *testing.T) {
	def := NewCard("Purifier", Sanctum, Creature, Rare, WithPower(5), WithArmor(2),
		WithConstantAbility(ConstantAbility{
			DisableTriggers: []Trigger{TriggerDestroyed},
		}))
	if got := constantText(&def); got != "Destroyed effects cannot trigger." {
		t.Errorf("constantText = %q, want %q", got, "Destroyed effects cannot trigger.")
	}
}
