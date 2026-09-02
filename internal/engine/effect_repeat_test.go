package engine

import "testing"

// TestRepeat covers running an effect once per count, with the choices made
// afresh each time, and the two ways a repetition can be misconfigured.
func TestRepeat(t *testing.T) {
	e := Repeat{
		Times: InPlay{Player: Controller, Type: Creature, House: Mars, Ready: true},
		Do:    DealDamage{Target: Target{Kind: TargetChosenCreature}, Amount: 2},
	}
	want := "for each friendly ready Mars creature, deal 2 damage to a creature"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if err := (Repeat{Do: e.Do}).validate(); err == nil {
		t.Error("a repetition with no count should be rejected")
	}
	if err := (Repeat{Times: e.Times}).validate(); err == nil {
		t.Error("a repetition with no effect should be rejected")
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("Martian A", Mars, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(NewCard("Martian B", Mars, Creature, Common, WithPower(3)), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 9), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// The default chooser takes the first candidate each time, so both
	// repetitions land on the same creature: 2 damage twice.
	atEnemy := Repeat{
		Times: e.Times,
		Do:    DealDamage{Target: Target{Kind: TargetChosenEnemyCreature}, Amount: 2},
	}
	atEnemy.Resolve(ctx)
	if g.Damage(enemy) != 4 {
		t.Errorf("damage = %d, want 4", g.Damage(enemy))
	}

	// No ready Mars creature on the other side, so nothing repeats.
	other := &EffectContext{Resolver: g, Controller: 1}
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
	ctx := &EffectContext{Resolver: g, Source: source, Controller: 0}

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
		Per: InPlay{
			Player: Controller,
			Type:   Creature,
			House:  Mars,
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
	ctx := &EffectContext{Resolver: g, Source: source, Controller: 0}

	e.Resolve(ctx)
	if g.Aember(1) != 4 {
		t.Errorf("opponent pool = %d, want 4 (the source does not count itself)", g.Aember(1))
	}
}
