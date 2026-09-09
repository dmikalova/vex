package engine

import "testing"

func TestRepeatOnCondition(t *testing.T) {
	e := RepeatOnCondition{
		Do:   Destroy{Target: Target{Kind: TargetEachEnemyCreature}.Selector(LeastPowerful)},
		Cond: Overwhelmed{},
	}
	if e.Text() != "destroy the least powerful enemy creature -> if you are overwhelmed, repeat this effect" {
		t.Errorf("text = %q", e.Text())
	}
	if (RepeatOnCondition{Do: Destroy{}, Cond: Overwhelmed{}}).validate() == nil {
		t.Error("unset destroy target should be invalid")
	}

	// Opponent has 3, controller 1: destroys until no longer overwhelmed.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 2), 1)
	g.AddToBattleline(testCreature("o2", 2), 1)
	g.AddToBattleline(testCreature("o3", 2), 1)
	g.AddToBattleline(testCreature("m1", 2), 0)
	RepeatOnCondition{
		Do:   Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Cond: Overwhelmed{},
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if len(g.Battleline(1)) != 1 {
		t.Errorf("opponent creatures = %d, want 1 (destroyed down to parity)", len(g.Battleline(1)))
	}

	// No enemy creatures: the gate stops the loop immediately.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(testCreature("m", 2), 0)
	RepeatOnCondition{
		Do:   Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Cond: Overwhelmed{},
	}.Resolve(
		&EffectContext{Resolver: g2, Controller: 0},
	)

	// A Do that cannot report progress always counts as progress, so only the
	// condition and the Rule of Six end the loop (Neutron Shark).
	g3 := NewGame("A", "B", 1)
	RepeatOnCondition{
		Do:   GainAember{Amount: 1, Player: Controller},
		Cond: PoolAember{Player: Opponent, Is: AtLeast, Amount: 1},
	}.Resolve(&EffectContext{Resolver: g3, Controller: 0})
	if g3.Aember(0) != 1 {
		t.Errorf("aember = %d, want 1 (the condition fails after one pass)", g3.Aember(0))
	}
}
