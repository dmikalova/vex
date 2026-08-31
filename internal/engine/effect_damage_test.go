package engine

import "testing"

func TestDamageIfDestroyed(t *testing.T) {
	t.Run("runs the follow-up only when the damage destroys the creature", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		weak := g.AddToBattleline(testCreature("weak", 2), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DamageIfDestroyed{
			Amount: 3,
			Target: Target{Kind: TargetChosenEnemyCreature},
			Then:   GainAember{Player: Controller, Amount: 1},
		}
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

		DamageIfDestroyed{
			Amount: 2,
			Target: Target{Kind: TargetChosenEnemyCreature},
			Then:   GainAember{Player: Controller, Amount: 1},
		}.Resolve(
			ctx,
		)
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

		DamageIfDestroyed{
			Amount: 3,
			Target: Target{Kind: TargetChosenEnemyCreature},
			Then:   PurgeCreature{Target: Target{Kind: TargetTriggeringCreature}},
		}.Resolve(
			ctx,
		)
		if g.State.Discard[1].contains(weak) {
			t.Error("the destroyed creature should be purged out of the discard pile")
		}
	})

	t.Run("no target is a no-op", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		DamageIfDestroyed{
			Amount: 3,
			Target: Target{Kind: TargetChosenEnemyCreature},
			Then:   GainAember{Player: Controller, Amount: 1},
		}.Resolve(
			ctx,
		)
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

	e := DealDamage{
		Amount: 1,
		Per:    InPlay{Player: Controller, Type: Creature},
		Target: Target{Kind: TargetEachEnemyCreature},
	}
	if e.Text() != "for each friendly creature in play, deal 1 damage to each enemy creature" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Damage(victim) != 3 {
		t.Errorf("victim damage = %d, want 3 (1 × 3 friendly creatures)", g.Damage(victim))
	}
}

func TestSpreadCreatureAndNeighbors(t *testing.T) {
	g := NewGame("A", "B", 1)
	// left and right are on flanks; only mid is a legal (non-flank) target.
	left := g.AddToBattleline(testCreature("left", 10), 1)
	mid := g.AddToBattleline(testCreature("mid", 10), 1)
	right := g.AddToBattleline(testCreature("right", 10), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DealDamage{Spread: CreatureAndNeighbors{Amount: 4, Splash: 2}}
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
	DealDamage{Spread: CreatureAndNeighbors{Amount: 4, Splash: 2}}.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Damage(a) != 0 {
		t.Errorf("no legal target should deal no damage, got %d", g2.Damage(a))
	}
}

func TestDealDamageIgnoreArmor(t *testing.T) {
	g := NewGame("A", "B", 1)
	armored := g.AddToBattleline(testCreature("armored", 5, WithArmor(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DealDamage{Amount: 3, Target: Target{Kind: TargetEachEnemyCreature}, IgnoreArmor: true}
	if e.Text() != "deal 3 damage to each enemy creature, ignoring armor" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Damage(armored) != 3 {
		t.Errorf("damage = %d, want 3 (armor ignored)", g.Damage(armored))
	}
}

func TestSpreadDifferentCreatures(t *testing.T) {
	t.Run("damages two different creatures", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 5), 1)
		b := g.AddToBattleline(testCreature("b", 5), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DealDamage{Spread: DifferentCreatures{First: 2, Second: 2}}
		if e.Text() != "deal 2 damage to a creature and deal 2 damage to a different creature" {
			t.Errorf("text = %q", e.Text())
		}
		e.Resolve(ctx)
		if g.Damage(a) != 2 || g.Damage(b) != 2 {
			t.Errorf("damage = %d/%d, want 2/2", g.Damage(a), g.Damage(b))
		}
	})

	t.Run("with only one creature, damages just it", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 5), 1)
		DealDamage{Spread: DifferentCreatures{First: 2, Second: 3}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(a) != 2 {
			t.Errorf("damage = %d, want 2", g.Damage(a))
		}
	})

	t.Run("with no creatures, does nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		DealDamage{Spread: DifferentCreatures{First: 2, Second: 2}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	})

	t.Run("validate rejects combining a Spread with a plain target", func(t *testing.T) {
		if (DealDamage{Spread: DifferentCreatures{First: 1, Second: 1}, Target: Target{Kind: TargetEachCreature}}).validate() == nil {
			t.Error("Spread + Target should be invalid")
		}
	})
}

func TestSpreadFlankWalk(t *testing.T) {
	t.Run("walks inward from the first flank dealing decreasing damage", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		left := g.AddToBattleline(testCreature("left", 9), 1)
		mid := g.AddToBattleline(testCreature("mid", 9), 1)
		right := g.AddToBattleline(testCreature("right", 9), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DealDamage{Spread: FlankWalk{Amounts: []int{3, 2, 1}}}
		if e.Text() != "choose a flank creature. Deal 3 damage to it, 2 damage to its neighbor, and 1 damage to the neighbor's other neighbor" {
			t.Errorf("text = %q", e.Text())
		}
		e.Resolve(ctx)
		if g.Damage(left) != 3 || g.Damage(mid) != 2 || g.Damage(right) != 1 {
			t.Errorf(
				"damage = %d/%d/%d, want 3/2/1",
				g.Damage(left),
				g.Damage(mid),
				g.Damage(right),
			)
		}
	})

	t.Run("walks the other direction from the far flank", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		left := g.AddToBattleline(testCreature("left", 9), 1)
		mid := g.AddToBattleline(testCreature("mid", 9), 1)
		right := g.AddToBattleline(testCreature("right", 9), 1)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{right}})
		DealDamage{Spread: FlankWalk{Amounts: []int{3, 2, 1}}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(right) != 3 || g.Damage(mid) != 2 || g.Damage(left) != 1 {
			t.Errorf(
				"damage = %d/%d/%d, want 3/2/1",
				g.Damage(right),
				g.Damage(mid),
				g.Damage(left),
			)
		}
	})

	t.Run("stops at the far flank when the line is shorter than the amounts", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 9), 1)
		b := g.AddToBattleline(testCreature("b", 9), 1)
		DealDamage{Spread: FlankWalk{Amounts: []int{3, 2, 1}}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(a) != 3 || g.Damage(b) != 2 {
			t.Errorf("damage = %d/%d, want 3/2", g.Damage(a), g.Damage(b))
		}
	})

	t.Run("with no flank creature, does nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		DealDamage{Spread: FlankWalk{Amounts: []int{3}}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	})

	t.Run("validate", func(t *testing.T) {
		if (DealDamage{Spread: FlankWalk{}}).validate() == nil {
			t.Error("empty amounts should be invalid")
		}
		if (DealDamage{Spread: FlankWalk{Amounts: []int{1}}}).validate() != nil {
			t.Error("a non-empty amounts list should be valid")
		}
	})
}
