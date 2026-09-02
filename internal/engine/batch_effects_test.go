package engine

import "testing"

func TestPreventDamage(t *testing.T) {
	g := NewGame("A", "B", 1)
	friend := g.AddToBattleline(testCreature("friend", 5), 0)
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PreventDamage{Target: Target{Kind: TargetEachFriendlyCreature}, Duration: EndOfTurn}
	if e.Text() != "for the remainder of the turn, each friendly creature cannot be dealt damage" {
		t.Errorf("text = %q", e.Text())
	}
	if (PreventDamage{Duration: EndOfTurn}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (PreventDamage{Target: Target{Kind: TargetEachFriendlyCreature}}).validate() == nil {
		t.Error("unset duration should be invalid")
	}
	if e.validate() != nil {
		t.Error("a set target and duration should be valid")
	}

	e.Resolve(ctx)
	g.applyRawDamage(friend, 3, false)
	if g.Damage(friend) != 0 {
		t.Errorf("protected creature took %d damage, want 0", g.Damage(friend))
	}

	// Protect an enemy creature too, then confirm end of turn clears both.
	PreventDamage{Target: Target{Kind: TargetEachEnemyCreature}, Duration: EndOfTurn}.Resolve(ctx)
	if !g.State.Cards[foe].DamageImmune {
		t.Fatal("enemy creature should be protected")
	}
	g.BeginTurn(0)
	g.EndTurn(0)
	if g.State.Cards[friend].DamageImmune || g.State.Cards[foe].DamageImmune {
		t.Error("end of turn should clear damage immunity on both players' creatures")
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

func TestHealGate(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 5), 0)
	g.State.Cards[c].Damage = 3
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if !(Heal{Fully: true, Target: Target{Kind: TargetChosenFriendlyCreature}}).resolveGate(ctx) {
		t.Error("healing a damaged creature should report a heal")
	}
	if !ctx.HasIt || ctx.It != c {
		t.Errorf("the healed creature should be left in context, got %d", ctx.It)
	}

	// Now undamaged: the gate reports nothing healed.
	ctx2 := &EffectContext{Resolver: g, Controller: 0}
	if (Heal{Fully: true, Target: Target{Kind: TargetChosenFriendlyCreature}}).resolveGate(ctx2) {
		t.Error("healing an undamaged creature should report no heal")
	}
}
