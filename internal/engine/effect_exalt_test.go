package engine

import "testing"

func TestExaltEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if got := (Exalt{Target: Target{Kind: TargetChosenFriendlyCreature}, Amount: 1}).Text(); got != "exalt a friendly creature" {
		t.Errorf("single exalt text = %q", got)
	}
	e := Exalt{Target: Target{Kind: TargetChosenEnemyCreature}, Amount: 2}
	if e.Text() != "exalt an enemy creature 2 times" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Cards[enemy].Amber != 2 {
		t.Errorf("amber on enemy = %d, want 2", g.State.Cards[enemy].Amber)
	}

	// No candidates: remove the enemy and resolve again (logs, no panic).
	g.DestroyEach(0, []LocalID{enemy})
	e.Resolve(ctx)
}
