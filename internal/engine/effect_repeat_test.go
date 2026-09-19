package engine

import "testing"

// TestRepeat covers running an effect once per count, with the choices made
// afresh each time, and the two ways a repetition can be misconfigured.
func TestRepeat(t *testing.T) {
	e := ForEach{
		Times: CardsInPlay{
			Player: Controller,
			Type:   Creature,
			House:  namedHouse(Mars),
			Ready:  true,
		},
		Do: DealDamage{
			Target: Target{Kind: TargetChosenCreature},
			Amount: 2,
		},
	}
	want := "for each friendly ready Mars creature, deal 2 damage to a creature"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if err := (ForEach{Do: e.Do}).validate(); err == nil {
		t.Error("a repetition with no count should be rejected")
	}
	if err := (ForEach{Times: e.Times}).validate(); err == nil {
		t.Error("a repetition with no effect should be rejected")
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("Martian A", Mars, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(NewCard("Martian B", Mars, Creature, Common, WithPower(3)), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 9), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	// The default chooser takes the first candidate each time, so both
	// repetitions land on the same creature: 2 damage twice.
	atEnemy := ForEach{
		Times: e.Times,
		Do: DealDamage{
			Target: Target{Kind: TargetChosenEnemyCreature},
			Amount: 2,
		},
	}
	atEnemy.Resolve(ctx)
	if g.Damage(enemy) != 4 {
		t.Errorf("damage = %d, want 4", g.Damage(enemy))
	}

	// No ready Mars creature on the other side, so nothing repeats.
	other := &EffectContext{
		Resolver:   g,
		Controller: 1,
	}
	atEnemy.Resolve(other)
	if g.Damage(enemy) != 4 {
		t.Errorf("damage after an empty repetition = %d, want 4", g.Damage(enemy))
	}
}

// TestCopiesInDiscard covers the count that pays a card off for having been
// played before.
func TestCopiesInDiscard(t *testing.T) {
	if got := (CopiesInDiscard{}).CountText(); got != "copy of "+SelfName+" in your discard pile" {
		t.Errorf("count text = %q", got)
	}

	g := NewGame("A", "B", 1)
	job := NewCard("Routine Job", Shadows, Tactic, Rare)
	source := g.AddToHand(job, 0)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     source,
		Controller: 0,
	}

	if got := (CopiesInDiscard{}).Value(ctx); got != 0 {
		t.Errorf("copies in an empty discard pile = %d, want 0", got)
	}

	g.AddToDiscard(job, 0)
	g.AddToDiscard(job, 0)
	g.AddToDiscard(NewCard("Other Job", Shadows, Tactic, Rare), 0)
	if got := (CopiesInDiscard{}).Value(ctx); got != 2 {
		t.Errorf("copies = %d, want 2", got)
	}
}

// TestLoseAemberPer covers scaling a fixed loss by a running count.
func TestLoseAemberPer(t *testing.T) {
	e := LoseAember{
		Player: Opponent,
		Amount: 1,
		Per: CardsInPlay{
			Player: Controller,
			Type:   Creature,
			House:  namedHouse(Mars),
			Other:  true,
		},
	}
	want := "for each other friendly Mars creature, your opponent loses 1 Æmber"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	source := g.AddToBattleline(NewCard("Phylyx", Mars, Creature, Rare, WithPower(1)), 0)
	g.AddToBattleline(NewCard("Martian", Mars, Creature, Common, WithPower(3)), 0)
	g.SetAember(1, 5)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     source,
		Controller: 0,
	}

	e.Resolve(ctx)
	if g.Aember(1) != 4 {
		t.Errorf("opponent pool = %d, want 4 (the source does not count itself)", g.Aember(1))
	}
}

func TestRepeatWhile(t *testing.T) {
	e := Repeat{
		Do:   Destroy{Target: Target{Kind: TargetEachEnemyCreature}.Refine(LeastPowerful)},
		Gate: While{Cond: Overwhelmed{}},
	}
	if e.Text() != "destroy the least powerful enemy creature -> if you are overwhelmed, repeat this effect" {
		t.Errorf("text = %q", e.Text())
	}
	if (Repeat{
		Do:   Destroy{},
		Gate: While{Cond: Overwhelmed{}},
	}).validate() == nil {
		t.Error("unset destroy target should be invalid")
	}

	// Opponent has 3, controller 1: destroys until no longer overwhelmed.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 2), 1)
	g.AddToBattleline(testCreature("o2", 2), 1)
	g.AddToBattleline(testCreature("o3", 2), 1)
	g.AddToBattleline(testCreature("m1", 2), 0)
	Repeat{
		Do:   Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Gate: While{Cond: Overwhelmed{}},
	}.Resolve(
		&EffectContext{
			Resolver:   g,
			Controller: 0,
		},
	)
	if len(g.Battleline(1)) != 1 {
		t.Errorf("opponent creatures = %d, want 1 (destroyed down to parity)", len(g.Battleline(1)))
	}

	// No enemy creatures: the gate stops the loop immediately.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(testCreature("m", 2), 0)
	Repeat{
		Do:   Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Gate: While{Cond: Overwhelmed{}},
	}.Resolve(
		&EffectContext{
			Resolver:   g2,
			Controller: 0,
		},
	)

	// A Do that cannot report progress always counts as progress, so only the
	// condition and the Rule of Six end the loop (Neutron Shark).
	g3 := NewGame("A", "B", 1)
	Repeat{
		Do: GainAember{
			Amount: 1,
			Player: Controller,
		},
		Gate: While{Cond: PoolAember{
			Player: Opponent,
			Is:     AtLeast,
			Amount: 1,
		}},
	}.Resolve(&EffectContext{
		Resolver:   g3,
		Controller: 0,
	})
	if g3.Aember(0) != 1 {
		t.Errorf("aember = %d, want 1 (the condition fails after one pass)", g3.Aember(0))
	}
}

// TestRepeatWhileSteal covers Bait and Switch: steal 1 Æmber, then repeat
// while the opponent still leads. A gating Do (the steal) that makes no progress
// ends the loop even while the condition holds.
func TestRepeatWhileSteal(t *testing.T) {
	e := Repeat{
		Do: StealAember{Amount: 1},
		Gate: While{Cond: PoolAember{
			Player: Opponent,
			Is:     MoreThanYou,
		}},
	}
	if e.Text() != "steal 1 Æmber -> if your opponent has more Æmber than you, repeat this effect" {
		t.Errorf("text = %q", e.Text())
	}

	// Opponent leads 5/0: steal until the lead is gone (5/0 -> 4/1 -> 3/2 -> 2/3).
	g := NewGame("A", "B", 1)
	g.State.Aember[0], g.State.Aember[1] = 0, 5
	e.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
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
	e.Resolve(&EffectContext{
		Resolver:   g2,
		Controller: 0,
	})
	if g2.Aember(0) != 0 || g2.Aember(1) != 5 {
		t.Errorf(
			"protected pool: you=%d opp=%d, want 0/5 (nothing stolen)",
			g2.Aember(0),
			g2.Aember(1),
		)
	}
}

// TestRepeatValidate covers node-level validation the gates delegate to.
func TestRepeatValidate(t *testing.T) {
	if (Repeat{Gate: While{Cond: Overwhelmed{}}}).validate() == nil {
		t.Error("missing Do should be invalid")
	}
	if (Repeat{Do: StealAember{Amount: 1}}).validate() == nil {
		t.Error("missing Gate should be invalid")
	}
	if (Repeat{
		Do:   StealAember{Amount: 1},
		Gate: ByExalting{},
	}).validate() == nil {
		t.Error("gate with unset target should be invalid")
	}
}
