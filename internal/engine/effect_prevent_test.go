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
	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.State.Cards[friend].DamageImmune || g.State.Cards[foe].DamageImmune {
		t.Error("end of turn should clear damage immunity on both players' creatures")
	}
}
