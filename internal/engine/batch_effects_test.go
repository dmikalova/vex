package engine

import "testing"

func TestPreventDamage(t *testing.T) {
	g := NewGame("A", "B", 1)
	friend := g.AddToBattleline(testCreature("friend", 5), 0)
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := PreventDamage{Target: Target{Kind: TargetEachFriendlyCreature}}
	if e.Text() != "for the remainder of the turn, each friendly creature cannot be dealt damage" {
		t.Errorf("text = %q", e.Text())
	}
	if (PreventDamage{}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if e.validate() != nil {
		t.Error("a set target should be valid")
	}

	e.Resolve(ctx)
	g.applyRawDamage(friend, 3, false)
	if g.Damage(friend) != 0 {
		t.Errorf("protected creature took %d damage, want 0", g.Damage(friend))
	}

	// Protect an enemy creature too, then confirm end of turn clears both.
	PreventDamage{Target: Target{Kind: TargetEachEnemyCreature}}.Resolve(ctx)
	if !g.State.Cards[foe].Invulnerable {
		t.Fatal("enemy creature should be protected")
	}
	g.BeginTurn(0)
	g.EndTurn(0)
	if g.State.Cards[friend].Invulnerable || g.State.Cards[foe].Invulnerable {
		t.Error("end of turn should clear damage immunity on both players' creatures")
	}
}

func TestMoveAemberToPool(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.AddAmberOn(c, 2)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if got := (MoveAemberToPool{}).Text(); got != "move 1 \u00c6mber from one of your cards to your pool" {
		t.Errorf("text = %q", got)
	}

	(MoveAemberToPool{}).Resolve(ctx)
	if g.AmberOn(c) != 1 || g.Aember(0) != 1 {
		t.Errorf("after move: card Æmber=%d pool=%d, want 1/1", g.AmberOn(c), g.Aember(0))
	}

	// No card carries Æmber: nothing happens.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(testCreature("bare", 3), 0)
	(MoveAemberToPool{}).Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Aember(0) != 0 {
		t.Error("moving with no Æmber-bearing card should do nothing")
	}
}

func TestMoveAemberToPoolDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	b := g.AddToBattleline(testCreature("b", 3), 0)
	g.AddAmberOn(a, 1)
	g.AddAmberOn(b, 1)
	g.SetChooser(0, orderRejectChooser{})
	(MoveAemberToPool{}).Resolve(&EffectContext{Resolver: g, Controller: 0})
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
