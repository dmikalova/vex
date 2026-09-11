package engine

import "testing"

func TestRepeatOnCondition(t *testing.T) {
	e := RepeatOnCondition{
		Do:   Destroy{Target: Target{Kind: TargetEachEnemyCreature}.Refine(LeastPowerful)},
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

// TestRepeatOnConditionSteal covers Bait and Switch: steal 1 Æmber, then repeat
// while the opponent still leads. A gating Do (the steal) that makes no progress
// ends the loop even while the condition holds.
func TestRepeatOnConditionSteal(t *testing.T) {
	e := RepeatOnCondition{
		Do:   StealAember{Amount: 1},
		Cond: PoolAember{Player: Opponent, Is: MoreThanYou},
	}
	if e.Text() != "steal 1 Æmber -> if your opponent has more Æmber than you, repeat this effect" {
		t.Errorf("text = %q", e.Text())
	}

	// Opponent leads 5/0: steal until the lead is gone (5/0 -> 4/1 -> 3/2 -> 2/3).
	g := NewGame("A", "B", 1)
	g.State.Aember[0], g.State.Aember[1] = 0, 5
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Aember(0) != 3 || g.Aember(1) != 2 {
		t.Errorf("after repeat: you=%d opp=%d, want 3/2", g.Aember(0), g.Aember(1))
	}

	// The opponent leads but their pool is protected, so the steal moves nothing.
	// The condition stays true, so the loop must stop on the action making no
	// progress rather than spin.
	g2 := NewGame("A", "B", 1)
	g2.State.Aember[0], g2.State.Aember[1] = 0, 5
	g2.AddToBattleline(
		NewCard("keeper", Sanctum, Creature, Rare, WithPower(4), WithAemberCannotBeStolen()),
		1,
	)
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Aember(0) != 0 || g2.Aember(1) != 5 {
		t.Errorf(
			"protected pool: you=%d opp=%d, want 0/5 (nothing stolen)",
			g2.Aember(0),
			g2.Aember(1),
		)
	}
}
