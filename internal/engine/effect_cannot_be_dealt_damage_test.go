package engine

import "testing"

func TestCannotBeDealtDamage(t *testing.T) {
	g := NewGame("A", "B", 1)
	friend := g.AddToBattleline(testCreature("friend", 5), 0)
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := CannotBeDealtDamage{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Duration: RemainderOfPlayerTurn,
	}
	if e.Text() != "for the remainder of the turn, each friendly creature cannot be dealt damage" {
		t.Errorf("text = %q", e.Text())
	}
	if (CannotBeDealtDamage{Duration: RemainderOfPlayerTurn}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (CannotBeDealtDamage{Target: Target{Kind: TargetEachFriendlyCreature}}).validate() == nil {
		t.Error("unset duration should be invalid")
	}
	if e.validate() != nil {
		t.Error("a set target and duration should be valid")
	}

	e.Resolve(ctx)
	g.applyRawDamage(DamageTarget{ID: friend, Amount: 3})
	if g.Damage(friend) != 0 {
		t.Errorf("protected creature took %d damage, want 0", g.Damage(friend))
	}
	// A friendly creature that arrives after the immunity resolves is protected too:
	// the side-wide mask is read live, not a snapshot of who was in play.
	late := g.AddToBattleline(testCreature("late", 5), 0)
	g.applyRawDamage(DamageTarget{ID: late, Amount: 3})
	if g.Damage(late) != 0 {
		t.Errorf("late-arriving friendly creature took %d damage, want 0", g.Damage(late))
	}

	// Protect the enemy side too, then confirm end of turn clears both.
	CannotBeDealtDamage{
		Target:   Target{Kind: TargetEachEnemyCreature},
		Duration: RemainderOfPlayerTurn,
	}.Resolve(
		ctx,
	)
	if !g.State.SideDamageImmune[1] {
		t.Fatal("enemy side should be protected")
	}
	g.applyRawDamage(DamageTarget{ID: foe, Amount: 3})
	if g.Damage(foe) != 0 {
		t.Errorf("protected enemy creature took %d damage, want 0", g.Damage(foe))
	}
	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.State.SideDamageImmune[0] || g.State.SideDamageImmune[1] {
		t.Error("end of turn should clear side-wide damage immunity for both players")
	}
}

// A per-card or filtered target is not a whole side, so it protects the concrete
// creatures it selects (and keeps rendering that phrase) rather than the side.
func TestCannotBeDealtDamageWholeSide(t *testing.T) {
	if _, ok := (Target{Kind: TargetThisCreature}).wholeSide(0); ok {
		t.Error("a single-creature target is not a whole side")
	}
	filtered := Target{Kind: TargetEachFriendlyCreature}.WithTrait(Knight)
	if _, ok := filtered.wholeSide(0); ok {
		t.Error("a trait-filtered friendly target is not a whole side")
	}
	if p, ok := (Target{Kind: TargetEachEnemyCreature}).wholeSide(0); !ok || p != 1 {
		t.Errorf("enemy side for controller 0 = (%d, %v), want (1, true)", p, ok)
	}
}
