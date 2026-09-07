package engine

import "testing"

func TestEnemyCreatureDestroyedReaction(t *testing.T) {
	if !EventEnemyCreatureDestroyed.isReaction() {
		t.Error("EventEnemyCreatureDestroyed should be a reaction point")
	}
	if EventEnemyCreatureDestroyed.clause() != "each time an enemy creature is destroyed" {
		t.Errorf("clause = %q", EventEnemyCreatureDestroyed.clause())
	}

	g := NewGame("A", "B", 1)
	g.AddLasting(
		LastingEffect{On: EventEnemyCreatureDestroyed, Do: actGainAember, Controller: 0, Amount: 1},
	)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	g.destroyEach(0, []LocalID{foe})
	if g.State.Aember[0] != 1 {
		t.Errorf(
			"controller Æmber = %d, want 1 after an enemy creature was destroyed",
			g.State.Aember[0],
		)
	}
}
