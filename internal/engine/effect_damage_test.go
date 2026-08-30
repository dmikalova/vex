package engine

import "testing"

func TestDamageIfDestroyed(t *testing.T) {
	t.Run("runs the follow-up only when the damage destroys the creature", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		weak := g.AddToBattleline(testCreature("weak", 2), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DamageIfDestroyed{Amount: 3, Target: Target{Kind: TargetChosenEnemyCreature}, Then: GainAember{Player: Controller, Amount: 1}}
		if e.Text() != "deal 3 damage to an enemy creature. If this damage destroys that creature, gain 1 Æmber" {
			t.Errorf("text = %q", e.Text())
		}
		e.Resolve(ctx)
		if resolverInPlay(ctx, weak) {
			t.Error("weak creature should be destroyed")
		}
		if g.State.Aember[0] != 1 {
			t.Errorf("aember = %d, want 1 (follow-up ran)", g.State.Aember[0])
		}
	})

	t.Run("does nothing extra when the creature survives", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		tough := g.AddToBattleline(testCreature("tough", 6), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		DamageIfDestroyed{Amount: 2, Target: Target{Kind: TargetChosenEnemyCreature}, Then: GainAember{Player: Controller, Amount: 1}}.Resolve(ctx)
		if !resolverInPlay(ctx, tough) {
			t.Error("tough creature should survive")
		}
		if g.State.Aember[0] != 0 {
			t.Errorf("aember = %d, want 0 (follow-up did not run)", g.State.Aember[0])
		}
	})

	t.Run("purges the destroyed creature via PurgeCreature", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		weak := g.AddToBattleline(testCreature("weak", 1), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		DamageIfDestroyed{Amount: 3, Target: Target{Kind: TargetChosenEnemyCreature}, Then: PurgeCreature{Target: Target{Kind: TargetTriggeringCreature}}}.Resolve(ctx)
		if g.State.Discard[1].contains(weak) {
			t.Error("the destroyed creature should be purged out of the discard pile")
		}
	})

	t.Run("no target is a no-op", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		DamageIfDestroyed{Amount: 3, Target: Target{Kind: TargetChosenEnemyCreature}, Then: GainAember{Player: Controller, Amount: 1}}.Resolve(ctx)
		if g.State.Aember[0] != 0 {
			t.Error("no target should not run the follow-up")
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (DamageIfDestroyed{Then: GainAember{Player: Controller, Amount: 1}}).validate() == nil {
			t.Error("unset target should be invalid")
		}
		if (DamageIfDestroyed{Target: Target{Kind: TargetChosenCreature}, Then: GainAember{}}).validate() == nil {
			t.Error("invalid follow-up should surface")
		}
	})
}

func TestDealDamagePerCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Three friendly creatures multiply the 1 base damage to 3.
	g.AddToBattleline(testCreature("f1", 5), 0)
	g.AddToBattleline(testCreature("f2", 5), 0)
	g.AddToBattleline(testCreature("f3", 5), 0)
	victim := g.AddToBattleline(testCreature("victim", 10), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DealDamage{Amount: 1, Per: InPlay{Player: Controller, Type: Creature}, Target: Target{Kind: TargetEachEnemyCreature}}
	if e.Text() != "for each friendly creature in play, deal 1 damage to each enemy creature" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Damage(victim) != 3 {
		t.Errorf("victim damage = %d, want 3 (1 × 3 friendly creatures)", g.Damage(victim))
	}
}

func TestSplashDamage(t *testing.T) {
	g := NewGame("A", "B", 1)
	// left and right are on flanks; only mid is a legal (non-flank) target.
	left := g.AddToBattleline(testCreature("left", 10), 1)
	mid := g.AddToBattleline(testCreature("mid", 10), 1)
	right := g.AddToBattleline(testCreature("right", 10), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := SplashDamage{Amount: 4, Splash: 2}
	if e.Text() != "deal 4 damage to a creature that is not on a flank and 2 damage to each of its neighbors" {
		t.Errorf("text = %q", e.Text())
	}
	// The default chooser picks the only non-flank creature, mid.
	e.Resolve(ctx)
	if g.Damage(mid) != 4 {
		t.Errorf("chosen damage = %d, want 4", g.Damage(mid))
	}
	if g.Damage(left) != 2 || g.Damage(right) != 2 {
		t.Errorf("neighbor damage = %d/%d, want 2/2", g.Damage(left), g.Damage(right))
	}

	// No non-flank creature (two creatures, both flanks): nothing happens.
	g2 := NewGame("A", "B", 1)
	a := g2.AddToBattleline(testCreature("a", 5), 1)
	g2.AddToBattleline(testCreature("b", 5), 1)
	SplashDamage{Amount: 4, Splash: 2}.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Damage(a) != 0 {
		t.Errorf("no legal target should deal no damage, got %d", g2.Damage(a))
	}
}
