package engine

import "testing"

// TestMoveAemberAll checks the All mode empties every source instead of moving a
// fixed amount, and that it cannot be combined with one.
func TestMoveAemberAll(t *testing.T) {
	e := MoveAember{
		All:  true,
		From: Target{Kind: TargetEachEnemyCreature},
		To:   Controller,
	}
	want := "move all Æmber from each enemy creature to your pool"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if (MoveAember{All: true, Amount: 1, From: e.From, To: Controller}).validate() == nil {
		t.Error("All combined with Amount should not validate")
	}

	g := NewGame("A", "B", 1)
	rich := g.AddToBattleline(testCreature("rich", 4), 1)
	poor := g.AddToBattleline(testCreature("poor", 4), 1)
	g.AddAmberOn(rich, 3)
	g.AddAmberOn(poor, 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.AmberOn(rich) != 0 || g.AmberOn(poor) != 0 {
		t.Errorf("Æmber left on the creatures = %d/%d, want 0/0",
			g.AmberOn(rich), g.AmberOn(poor))
	}
	if got := g.Aember(0); got != 4 {
		t.Errorf("pool = %d, want 4", got)
	}
}

func TestMoveAember(t *testing.T) {
	friendly := Target{Kind: TargetChosenFriendlyCreatureOrArtifact}

	// Text: pool destination and card destination.
	toPool := MoveAember{Amount: 1, From: friendly, To: Controller}
	if got := toPool.Text(); got != "move 1 \u00c6mber from a friendly creature or artifact to your pool" {
		t.Errorf("pool text = %q", got)
	}
	toOpp := MoveAember{Amount: 2, From: friendly, To: Opponent}
	if got := toOpp.Text(); got != "move 2 \u00c6mber from a friendly creature or artifact to your opponent's pool" {
		t.Errorf("opponent-pool text = %q", got)
	}
	toCard := MoveAember{From: friendly, Onto: Target{Kind: TargetChosenEnemyCreature}}
	if got := toCard.Text(); got != "move 1 \u00c6mber from a friendly creature or artifact to an enemy creature" {
		t.Errorf("card text = %q", got)
	}

	// validate: source and exactly one destination.
	if (MoveAember{To: Controller}).validate() == nil {
		t.Error("unset source should be invalid")
	}
	if (MoveAember{From: friendly}).validate() == nil {
		t.Error("no destination should be invalid")
	}
	if (MoveAember{From: friendly, To: Controller, Onto: friendly}).validate() == nil {
		t.Error("two destinations should be invalid")
	}
	if toPool.validate() != nil || toCard.validate() != nil {
		t.Error("one destination should be valid")
	}

	// Resolve into a pool, capping the move at what the source holds.
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.AddAmberOn(c, 2)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	MoveAember{Amount: 3, From: friendly, To: Controller}.Resolve(ctx)
	if g.AmberOn(c) != 0 || g.Aember(0) != 2 {
		t.Errorf("after capped move: card=%d pool=%d, want 0/2", g.AmberOn(c), g.Aember(0))
	}

	// Resolve onto another card.
	g2 := NewGame("A", "B", 1)
	src := g2.AddToBattleline(testCreature("src", 3), 0)
	dst := g2.AddToBattleline(testCreature("dst", 3), 1)
	g2.AddAmberOn(src, 2)
	MoveAember{Amount: 1, From: friendly, Onto: Target{Kind: TargetChosenEnemyCreature}}.
		Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.AmberOn(src) != 1 || g2.AmberOn(dst) != 1 {
		t.Errorf("after card move: src=%d dst=%d, want 1/1", g2.AmberOn(src), g2.AmberOn(dst))
	}

	// No source carries Æmber: nothing happens.
	g3 := NewGame("A", "B", 1)
	g3.AddToBattleline(testCreature("bare", 3), 0)
	MoveAember{From: friendly, To: Controller}.Resolve(&EffectContext{Resolver: g3, Controller: 0})
	if g3.Aember(0) != 0 {
		t.Error("moving with no Æmber-bearing card should do nothing")
	}

	// A card destination with no candidate card: nothing moves.
	g4 := NewGame("A", "B", 1)
	only := g4.AddToBattleline(testCreature("only", 3), 0)
	g4.AddAmberOn(only, 1)
	MoveAember{From: friendly, Onto: Target{Kind: TargetChosenEnemyCreature}}.
		Resolve(&EffectContext{Resolver: g4, Controller: 0})
	if g4.AmberOn(only) != 1 {
		t.Error("a move with no destination card should move nothing")
	}
}

func TestMoveAemberDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	b := g.AddToBattleline(testCreature("b", 3), 0)
	g.AddAmberOn(a, 1)
	g.AddAmberOn(b, 1)
	g.SetChooser(0, orderRejectChooser{})
	MoveAember{From: Target{Kind: TargetChosenFriendlyCreatureOrArtifact}, To: Controller}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Aember(0) != 0 {
		t.Error("a declined move should move nothing")
	}
}
