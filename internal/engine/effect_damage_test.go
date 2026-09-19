package engine

import "testing"

// A "you may deal damage to a creature" is one card choice, so the player picks
// the creature directly instead of first answering Yes (Rock-Hurling Giant).
func TestMayDeclinableDealDamage(t *testing.T) {
	e := May{Do: DealDamage{Amount: 4, Target: Target{Kind: TargetChosenCreature}}}
	if !e.Do.(declinableEffect).declinable() {
		t.Fatal("a chosen-target DealDamage should be declinable")
	}

	g := NewGame("A", "B", 1)
	ch := &cardDecliner{}
	g.SetChooser(0, ch)
	victim := g.AddToBattleline(testCreature("Victim", 5), 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if ch.asked != 1 {
		t.Errorf("declinable prompts = %d, want 1", ch.asked)
	}
	if g.Damage(victim) != 4 {
		t.Errorf("the chosen creature should have taken 4 damage, got %d", g.Damage(victim))
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	other := declined.AddToBattleline(testCreature("Victim", 5), 1)
	e.Resolve(&EffectContext{Resolver: declined, Controller: 0})
	if declined.Damage(other) != 0 {
		t.Error("a declined May should deal no damage")
	}

	// A computed amount of zero is nothing to offer, so the candidate is never
	// asked for.
	if (DealDamage{Amount: 0, Target: Target{Kind: TargetChosenCreature}}).resolveOptional(
		&EffectContext{Resolver: g, Controller: 0},
	) {
		t.Error("a zero amount should resolve to nothing")
	}
}

func TestDamageThenIfDestroyed(t *testing.T) {
	t.Run("runs the follow-up only when the damage destroys the creature", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		weak := g.AddToBattleline(testCreature("weak", 2), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DealDamage{
			Amount: 3,
			After:  IfDestroyed,
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

		DealDamage{
			Amount: 2,
			After:  IfDestroyed,
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

		DealDamage{
			Amount: 3,
			After:  IfDestroyed,
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
		DealDamage{
			Amount: 3,
			After:  IfDestroyed,
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
		if (DealDamage{After: IfDestroyed, Then: GainAember{Player: Controller, Amount: 1}}).validate() == nil {
			t.Error("unset target should be invalid")
		}
		if (DealDamage{Target: Target{Kind: TargetChosenCreature}, Then: GainAember{Player: Controller, Amount: 1}}).validate() == nil {
			t.Error("unset After should be invalid")
		}
		if (DealDamage{Target: Target{Kind: TargetChosenCreature}, After: IfDestroyed, Then: GainAember{}}).validate() == nil {
			t.Error("invalid follow-up should surface")
		}
	})

	t.Run("follow-up can hit the destroyed creature's former neighbors", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		left := g.AddToBattleline(testCreature("left", 6), 1)
		mid := g.AddToBattleline(testCreature("mid", 2), 1)
		right := g.AddToBattleline(testCreature("right", 6), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DealDamage{
			Amount: 3,
			After:  IfDestroyed,
			Target: Target{Kind: TargetChosenEnemyCreature}.PowerAtMost(2),
			Then:   DealDamage{Amount: 2, Target: Target{Kind: TargetFormerNeighbors}},
		}
		if got := (Target{Kind: TargetFormerNeighbors}).Text(); got != "each of that creature's neighbors" {
			t.Errorf("former-neighbors text = %q", got)
		}
		e.Resolve(ctx)
		if resolverInPlay(ctx, mid) {
			t.Error("mid creature should be destroyed")
		}
		if g.Damage(left) != 2 || g.Damage(right) != 2 {
			t.Errorf("neighbor damage = %d/%d, want 2/2", g.Damage(left), g.Damage(right))
		}
	})
}

func TestDamageThen(t *testing.T) {
	t.Run("runs the follow-up whether or not the creature is destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		tough := g.AddToBattleline(testCreature("tough", 6), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DealDamage{
			Amount: 2,
			After:  Always,
			Target: Target{Kind: TargetChosenCreature},
			Then:   GainAember{Player: Controller, Amount: 1},
		}
		if got, want := e.Text(), "deal 2 damage to a creature and gain 1 \u00c6mber"; got != want {
			t.Fatalf("text = %q, want %q", got, want)
		}
		e.Resolve(ctx)
		if g.Damage(tough) != 2 {
			t.Errorf("damage = %d, want 2", g.Damage(tough))
		}
		if g.State.Aember[0] != 1 {
			t.Errorf("aember = %d, want 1 (follow-up ran on the survivor)", g.State.Aember[0])
		}
	})

	t.Run("puts the damaged creature in context", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		weak := g.AddToBattleline(testCreature("weak", 1), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		DealDamage{
			Amount: 3,
			After:  Always,
			Target: Target{Kind: TargetChosenEnemyCreature},
			Then:   PurgeCreature{Target: Target{Kind: TargetTriggeringCreature}},
		}.Resolve(ctx)
		if g.State.Discard[1].contains(weak) {
			t.Error("the destroyed creature should be purged out of the discard pile")
		}
	})

	t.Run("no target is a no-op", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		DealDamage{
			Amount: 3,
			After:  Always,
			Target: Target{Kind: TargetChosenEnemyCreature},
			Then:   GainAember{Player: Controller, Amount: 1},
		}.Resolve(ctx)
		if g.State.Aember[0] != 0 {
			t.Error("no target should not run the follow-up")
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (DealDamage{After: Always, Then: GainAember{Player: Controller, Amount: 1}}).validate() == nil {
			t.Error("unset target should be invalid")
		}
		if (DealDamage{Target: Target{Kind: TargetChosenCreature}, After: Always, Then: GainAember{}}).validate() == nil {
			t.Error("invalid follow-up should surface")
		}
	})
}

// TestForEachHouseDealsDamagePerHouse covers Gleeful Mayhem: for each house, one
// chosen creature of that house takes the damage, houses with no creature are
// skipped, and a second creature of an already-hit house is left untouched.
func TestForEachHouseDealsDamagePerHouse(t *testing.T) {
	if err := (ForEachHouse{}).validate(); err == nil {
		t.Error("want validate error for missing Do")
	}
	e := ForEachHouse{Do: DealDamage{
		Amount: 5,
		Target: Target{Kind: TargetChosenCreature}.House(HouseMatcher{Kind: MatchEachHouse}),
	}}
	if got := e.Text(); got != "for each house, deal 5 damage to a creature of that house" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	dis := g.AddToBattleline(NewCard("dis", Dis, Creature, Common, WithPower(10)), 0)
	otherDis := g.AddToBattleline(NewCard("dis2", Dis, Creature, Common, WithPower(10)), 0)
	logos := g.AddToBattleline(NewCard("logos", Logos, Creature, Common, WithPower(10)), 1)
	// The default chooser takes the first candidate for each house.
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e.Resolve(ctx)
	if g.Damage(dis) != 5 {
		t.Errorf("dis damage = %d, want 5", g.Damage(dis))
	}
	if g.Damage(otherDis) != 0 {
		t.Errorf("otherDis damage = %d, want 0 (only one per house)", g.Damage(otherDis))
	}
	if g.Damage(logos) != 5 {
		t.Errorf("logos damage = %d, want 5", g.Damage(logos))
	}

	if err := e.validate(); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}

	// Declining the choice for a house deals no damage to that house.
	g2 := NewGame("A", "B", 1)
	safe := g2.AddToBattleline(NewCard("safe", Dis, Creature, Common, WithPower(10)), 0)
	g2.AddToBattleline(NewCard("safe2", Dis, Creature, Common, WithPower(10)), 0)
	g2.SetChooser(0, orderRejectChooser{})
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Damage(safe) != 0 {
		t.Error("a declined house should take no damage")
	}
}

func TestDealDamagePerCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Three friendly creatures multiply the 1 base damage to 3.
	g.AddToBattleline(testCreature("f1", 5), 0)
	g.AddToBattleline(testCreature("f2", 5), 0)
	g.AddToBattleline(testCreature("f3", 5), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 10), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DealDamage{
		Amount: 1,
		Per:    CardsInPlay{Player: Controller, Type: Creature},
		Target: Target{Kind: TargetEachEnemyCreature},
	}
	if e.Text() != "for each friendly creature in play, deal 1 damage to each enemy creature" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Damage(enemy) != 3 {
		t.Errorf("enemy damage = %d, want 3 (1 × 3 friendly creatures)", g.Damage(enemy))
	}
}

// A "for each" DealDamage aimed at a single chosen creature makes one target
// choice per instance — each free to name a different creature, a creature named
// twice taking both hits at once — rather than one target taking the whole scaled
// amount (which is Sack of Coins' ChooseCreatureThen shape instead).
func TestDealDamagePerInstanceChoosesEachTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Three friendly creatures set the count to three instances of 1 damage.
	g.AddToBattleline(testCreature("f1", 5), 0)
	g.AddToBattleline(testCreature("f2", 5), 0)
	g.AddToBattleline(testCreature("f3", 5), 0)
	e1 := g.AddToBattleline(testCreature("e1", 10), 1)
	e2 := g.AddToBattleline(testCreature("e2", 10), 1)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{e1, e1, e2}})
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DealDamage{
		Amount: 1,
		Per:    CardsInPlay{Player: Controller, Type: Creature},
		Target: Target{Kind: TargetChosenEnemyCreature},
	}
	e.Resolve(ctx)
	if g.Damage(e1) != 2 {
		t.Errorf("e1 damage = %d, want 2 (chosen for two instances)", g.Damage(e1))
	}
	if g.Damage(e2) != 1 {
		t.Errorf("e2 damage = %d, want 1", g.Damage(e2))
	}
}

func TestDealDamagePerInstanceDegenerate(t *testing.T) {
	base := DealDamage{
		Amount: 1,
		Per:    CardsInPlay{Player: Controller, Type: Creature},
		Target: Target{Kind: TargetChosenEnemyCreature},
	}

	t.Run("a zero count asks nothing and deals nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToBattleline(testCreature("e1", 10), 1)
		ch := &countingChooser{}
		g.SetChooser(0, ch)
		base.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if ch.calls != 0 {
			t.Errorf("chooser asked %d times, want 0 (no friendly creatures to count)", ch.calls)
		}
	})

	t.Run("no candidate asks nothing and deals nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToBattleline(testCreature("f1", 5), 0) // count is one, but there is no enemy
		ch := &countingChooser{}
		g.SetChooser(0, ch)
		base.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if ch.calls != 0 {
			t.Errorf("chooser asked %d times, want 0 (no enemy creature to choose)", ch.calls)
		}
	})

	t.Run("declining every instance deals nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToBattleline(testCreature("f1", 5), 0)
		g.AddToBattleline(testCreature("f2", 5), 0) // count is two instances
		e1 := g.AddToBattleline(testCreature("e1", 10), 1)
		e2 := g.AddToBattleline(
			testCreature("e2", 10),
			1,
		) // two candidates, so the chooser is asked
		g.SetChooser(0, orderRejectChooser{}) // always declines
		base.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(e1) != 0 || g.Damage(e2) != 0 {
			t.Errorf("damage e1=%d e2=%d, want 0 and 0 when every choice is declined",
				g.Damage(e1), g.Damage(e2))
		}
	})
}

func TestSpreadCreatureAndNeighbors(t *testing.T) {
	// Without NotOnFlank, any creature is a legal target.
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 10), 1)
	mid := g.AddToBattleline(testCreature("mid", 10), 1)
	right := g.AddToBattleline(testCreature("right", 10), 1)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{mid}})
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DealDamage{Spread: CreatureAndNeighbors{Amount: 4, Splash: 2}}
	if e.Text() != "choose a creature. Deal 4 damage to the chosen creature and 2 damage to each of its neighbors" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Damage(mid) != 4 {
		t.Errorf("chosen damage = %d, want 4", g.Damage(mid))
	}
	if g.Damage(left) != 2 || g.Damage(right) != 2 {
		t.Errorf("neighbor damage = %d/%d, want 2/2", g.Damage(left), g.Damage(right))
	}

	// A flank creature is now a legal target and has a single neighbor.
	g2 := NewGame("A", "B", 1)
	flank := g2.AddToBattleline(testCreature("flank", 10), 1)
	inner := g2.AddToBattleline(testCreature("inner", 10), 1)
	DealDamage{Spread: CreatureAndNeighbors{Amount: 4, Splash: 2}}.Resolve(
		&EffectContext{Resolver: g2, Controller: 0},
	)
	if g2.Damage(flank) != 4 || g2.Damage(inner) != 2 {
		t.Errorf("flank hit = %d, neighbor = %d, want 4/2", g2.Damage(flank), g2.Damage(inner))
	}

	// NotOnFlank keeps the choice to an interior creature, which the default
	// chooser then picks as the only legal target.
	g3 := NewGame("A", "B", 1)
	fl := g3.AddToBattleline(testCreature("fl", 10), 1)
	center := g3.AddToBattleline(testCreature("center", 10), 1)
	fr := g3.AddToBattleline(testCreature("fr", 10), 1)
	restricted := DealDamage{Spread: CreatureAndNeighbors{Amount: 4, Splash: 2, NotOnFlank: true}}
	if restricted.Text() != "choose a creature that is not on a flank. "+
		"Deal 4 damage to the chosen creature and 2 damage to each of its neighbors" {
		t.Errorf("restricted text = %q", restricted.Text())
	}
	restricted.Resolve(&EffectContext{Resolver: g3, Controller: 0})
	if g3.Damage(center) != 4 || g3.Damage(fl) != 2 || g3.Damage(fr) != 2 {
		t.Errorf("restricted damage = %d (%d/%d), want 4 (2/2)",
			g3.Damage(center), g3.Damage(fl), g3.Damage(fr))
	}

	// NotOnFlank with only flank creatures: no legal target, nothing happens.
	g4 := NewGame("A", "B", 1)
	a := g4.AddToBattleline(testCreature("a", 5), 1)
	g4.AddToBattleline(testCreature("b", 5), 1)
	DealDamage{Spread: CreatureAndNeighbors{Amount: 4, Splash: 2, NotOnFlank: true}}.Resolve(
		&EffectContext{Resolver: g4, Controller: 0},
	)
	if g4.Damage(a) != 0 {
		t.Errorf("no legal target should deal no damage, got %d", g4.Damage(a))
	}
}

// TestSpreadCreatureAndNeighborsAtTarget proves that a named Target aims the
// spread at that creature — the one this creature fought (ctx.It) — hitting it
// and its neighbors with no choose-creature prompt. A chooser primed to pick a
// decoy would land the damage there instead if a choice were made.
func TestSpreadCreatureAndNeighborsAtTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 10), 1)
	fought := g.AddToBattleline(testCreature("fought", 10), 1)
	right := g.AddToBattleline(testCreature("right", 10), 1)
	decoy := g.AddToBattleline(testCreature("decoy", 10), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{decoy}})
	ctx := &EffectContext{Resolver: g, Controller: 0, It: fought, HasIt: true}

	e := DealDamage{Spread: CreatureAndNeighbors{
		Amount: 2,
		Splash: 2,
		Target: Target{Kind: TargetCreatureFought},
	}}
	want := "deal 2 damage to the creature {self} fought and each of its neighbors"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	e.Resolve(ctx)
	if g.Damage(fought) != 2 {
		t.Errorf("fought damage = %d, want 2", g.Damage(fought))
	}
	if g.Damage(left) != 2 || g.Damage(right) != 2 {
		t.Errorf("neighbor damage = %d/%d, want 2/2", g.Damage(left), g.Damage(right))
	}
	if g.Damage(decoy) != 0 {
		t.Errorf("decoy took %d damage; target should not be chosen", g.Damage(decoy))
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
		if e.Text() != "deal 2 damage to a creature and 2 damage to a different creature" {
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
		DealDamage{
			Spread: DifferentCreatures{First: 2, Second: 3},
		}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
		if g.Damage(a) != 2 {
			t.Errorf("damage = %d, want 2", g.Damage(a))
		}
	})

	t.Run("with no creatures, does nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		DealDamage{
			Spread: DifferentCreatures{First: 2, Second: 2},
		}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
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
		DealDamage{
			Spread: FlankWalk{Amounts: []int{3, 2, 1}},
		}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
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
		DealDamage{
			Spread: FlankWalk{Amounts: []int{3, 2, 1}},
		}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
		if g.Damage(a) != 3 || g.Damage(b) != 2 {
			t.Errorf("damage = %d/%d, want 3/2", g.Damage(a), g.Damage(b))
		}
	})

	t.Run("with no flank creature, does nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		DealDamage{
			Spread: FlankWalk{Amounts: []int{3}},
		}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	})

	t.Run("validate", func(t *testing.T) {
		if (DealDamage{Spread: FlankWalk{}}).validate() == nil {
			t.Error("empty amounts should be invalid")
		}
		if (DealDamage{Spread: FlankWalk{Amounts: []int{1}}}).validate() != nil {
			t.Error("a non-empty amounts list should be valid")
		}
		if (DealDamage{Spread: CreatureAndNeighbors{}}).validate() != nil {
			t.Error("a spread that cannot be misconfigured should be valid")
		}
	})
}

// declineAfterChooser answers with each queued id in turn, then declines.
type declineAfterChooser struct{ ids []LocalID }

func (c *declineAfterChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	return cands[0], true
}

func (c *declineAfterChooser) ChooseCardOrDecline(_, _ string, _ []LocalID) (LocalID, bool) {
	if len(c.ids) > 0 {
		id := c.ids[0]
		c.ids = c.ids[1:]
		return id, true
	}
	return 0, false
}

func TestSpreadUpToCreatures(t *testing.T) {
	t.Run("damages up to Creatures creatures and tallies the kills", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToBattleline(testCreature("a", 1), 1)
		g.AddToBattleline(testCreature("b", 1), 1)
		g.AddToBattleline(testCreature("c", 1), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DealDamage{Spread: UpToCreatures{Creatures: 3, Amount: 1}}
		if e.Text() != "deal 1 damage to up to 3 creatures" {
			t.Errorf("text = %q", e.Text())
		}
		e.Resolve(ctx)
		if got := ctx.Produced.TotalDestroyed(); got != 3 {
			t.Errorf("destroyed tally = %d, want 3", got)
		}
	})

	t.Run("survivors are not tallied", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 5), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		DealDamage{Spread: UpToCreatures{Creatures: 3, Amount: 1}}.Resolve(ctx)
		if g.Damage(a) != 1 {
			t.Errorf("damage = %d, want 1", g.Damage(a))
		}
		if got := ctx.Produced.TotalDestroyed(); got != 0 {
			t.Errorf("destroyed tally = %d, want 0", got)
		}
	})

	t.Run("stops when the controller declines", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 5), 1)
		b := g.AddToBattleline(testCreature("b", 5), 1)
		g.SetChooser(0, &declineAfterChooser{ids: []LocalID{a}})
		ctx := &EffectContext{Resolver: g, Controller: 0}
		DealDamage{Spread: UpToCreatures{Creatures: 3, Amount: 1}}.Resolve(ctx)
		if g.Damage(a) != 1 || g.Damage(b) != 0 {
			t.Errorf("damage = %d/%d, want 1/0", g.Damage(a), g.Damage(b))
		}
	})

	t.Run("with no creatures, does nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		DealDamage{Spread: UpToCreatures{Creatures: 3, Amount: 1}}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	})

	t.Run("validate rejects a Count below one", func(t *testing.T) {
		if (DealDamage{Spread: UpToCreatures{Creatures: 0, Amount: 1}}).validate() == nil {
			t.Error("Count 0 should be invalid")
		}
		if (DealDamage{Spread: UpToCreatures{Creatures: 3, Amount: 1}}).validate() != nil {
			t.Error("Count 3 should be valid")
		}
	})

	t.Run("Undamaged offers only creatures with no damage on them", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		hurt := g.AddToBattleline(testCreature("hurt", 9), 1)
		clean := g.AddToBattleline(testCreature("clean", 9), 1)
		g.SetDamage(hurt, 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DealDamage{Spread: UpToCreatures{Creatures: 3, Amount: 2, Undamaged: true}}
		if got := e.Text(); got != "deal 2 damage to up to 3 undamaged creatures" {
			t.Errorf("text = %q", got)
		}
		e.Resolve(ctx)
		if got := g.Damage(clean); got != 2 {
			t.Errorf("undamaged creature = %d, want 2", got)
		}
		if got := g.Damage(hurt); got != 1 {
			t.Errorf("already-damaged creature = %d, want 1 (untouched)", got)
		}
	})

	t.Run("WhenDamaged deals the larger amount to already-damaged creatures", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		hurt := g.AddToBattleline(testCreature("hurt", 9), 1)
		clean := g.AddToBattleline(testCreature("clean", 9), 1)
		g.SetDamage(hurt, 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := DealDamage{Spread: UpToCreatures{Creatures: 3, Amount: 1, WhenDamaged: 3}}
		if got := e.Text(); got !=
			"choose up to 3 creatures. Deal 1 damage to each chosen creature. "+
				"Deal 3 damage instead to each chosen creature that was already damaged" {
			t.Errorf("text = %q", got)
		}
		e.Resolve(ctx)
		if got := g.Damage(hurt); got != 4 {
			t.Errorf("already-damaged creature = %d, want 4 (1+3)", got)
		}
		if got := g.Damage(clean); got != 1 {
			t.Errorf("undamaged creature = %d, want 1", got)
		}
	})

	t.Run("validate rejects WhenDamaged combined with Undamaged", func(t *testing.T) {
		if (DealDamage{Spread: UpToCreatures{Creatures: 2, Amount: 1, WhenDamaged: 3, Undamaged: true}}).
			validate() == nil {
			t.Error("WhenDamaged + Undamaged should be invalid")
		}
	})
}

// TestDealDamagePerTarget checks the damage is scaled by a quantity read off
// each creature separately, so one batch can hit for different amounts.
func TestDealDamagePerTarget(t *testing.T) {
	e := DealDamage{
		Amount:    1,
		Target:    Target{Kind: TargetEachEnemyCreature},
		PerTarget: AemberOnIt,
	}
	want := "deal 1 damage to each enemy creature for each Æmber on it"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	withPer := e
	withPer.Per = CardsDestroyed{}
	if withPer.validate() == nil {
		t.Error("PerTarget combined with Per should not validate")
	}

	g := NewGame("A", "B", 1)
	laden := g.AddToBattleline(testCreature("laden", 9), 1)
	bare := g.AddToBattleline(testCreature("bare", 9), 1)
	g.AddAmberOn(laden, 2)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := g.Damage(laden); got != 2 {
		t.Errorf("damage on the laden creature = %d, want 2", got)
	}
	if got := g.Damage(bare); got != 0 {
		t.Errorf("damage on the bare creature = %d, want 0", got)
	}
}

// TestDealDamageDamageOnIt checks the DamageOnIt PerTarget scales each hit by the
// damage already sitting on that creature (Cauldron Boil).
func TestDealDamageDamageOnIt(t *testing.T) {
	e := DealDamage{
		Amount:    1,
		Target:    Target{Kind: TargetEachCreature},
		PerTarget: DamageOnIt,
	}
	want := "deal 1 damage to each creature for each point of damage on it"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	g := NewGame("A", "B", 1)
	hurt := g.AddToBattleline(testCreature("hurt", 9), 1)
	fine := g.AddToBattleline(testCreature("fine", 9), 1)
	g.SetDamage(hurt, 3)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	// The wounded creature takes 1 damage per point already on it (3 * 1 = 3),
	// landing on 6; the unwounded one takes nothing.
	if got := g.Damage(hurt); got != 6 {
		t.Errorf("damage on the wounded creature = %d, want 6", got)
	}
	if got := g.Damage(fine); got != 0 {
		t.Errorf("damage on the healthy creature = %d, want 0", got)
	}
}

// damageChooserTrap fails the test if the engine asks it anything: dealing zero
// damage must not reach a prompt at all.
type damageChooserTrap struct{ t *testing.T }

func (c damageChooserTrap) ChooseCreature(_, prompt string, _ []LocalID) (LocalID, bool) {
	c.t.Errorf("unexpected prompt %q: dealing 0 damage should ask nothing", prompt)
	return 0, false
}

// TestDealDamageZeroAmountAsksNothing pins that a computed amount of zero skips
// the effect entirely rather than prompting for a target it cannot hurt — the
// second half of Guardian Demon when its heal removed no damage.
func TestDealDamageZeroAmountAsksNothing(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("e1", 5), 1)
	g.AddToBattleline(testCreature("e2", 5), 1)
	g.SetChooser(0, damageChooserTrap{t})
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Produced.DamageHealed is 0: nothing was healed, so nothing is dealt.
	DealDamage{
		AmountFrom: DamageHealed{},
		Target:     Target{Kind: TargetChosenEnemyCreature},
	}.Resolve(ctx)
}

// midThenFirst picks mid for the first creature choice, then the first candidate
// for any later choice (the neighbor pick).
type midThenFirst struct {
	FirstChooser
	mid LocalID
}

func (c midThenFirst) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	for _, x := range cands {
		if x == c.mid {
			return x, true
		}
	}
	return cands[0], true
}

func TestDamageThenIfSurvives(t *testing.T) {
	g := NewGame("A", "B", 1)
	survivor := g.AddToBattleline(testCreature("s", 5), 1)
	g.AddToHand(testCreature("h1", 1), 1)
	g.AddToHand(testCreature("h2", 1), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DealDamage{
		Amount: 3,
		After:  IfSurvives,
		Target: Target{Kind: TargetChosenCreature},
		Then:   DiscardCard{Player: ItsOwner, Zones: []Zone{Hand}, Selection: Random{}},
	}
	if e.Text() != "deal 3 damage to a creature. If it is not destroyed, its owner discards a random card from their hand" {
		t.Errorf("text = %q", e.Text())
	}
	if (DealDamage{After: IfSurvives, Then: DiscardCard{Player: Opponent, Zones: []Zone{Hand}, Selection: Random{}}}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (DealDamage{Target: Target{Kind: TargetChosenCreature}, After: IfSurvives, Then: DiscardCard{Player: Opponent, Zones: []Zone{Hand}, Selection: Random{}}}).validate() != nil {
		t.Error("a set target with a valid follow-up should pass")
	}

	e.Resolve(ctx)
	if g.Damage(survivor) != 3 {
		t.Errorf("survivor damage = %d, want 3", g.Damage(survivor))
	}
	if g.State.Hand[1].Count != 1 {
		t.Errorf("owner hand = %d, want 1 (a card discarded)", g.State.Hand[1].Count)
	}

	// A creature that dies triggers no follow-up.
	g2 := NewGame("A", "B", 1)
	dead := g2.AddToBattleline(testCreature("d", 2), 1)
	g2.AddToHand(testCreature("keep", 1), 1)
	(DealDamage{Amount: 3, After: IfSurvives, Target: Target{Kind: TargetChosenCreature}, Then: DiscardCard{Player: ItsOwner, Zones: []Zone{Hand}, Selection: Random{}}}).Resolve(
		&EffectContext{Resolver: g2, Controller: 0},
	)
	if g2.inPlay(dead) {
		t.Error("the creature should have been destroyed")
	}
	if g2.State.Hand[1].Count != 1 {
		t.Error("a destroyed creature's owner should not discard")
	}

	// No creature to target: nothing happens.
	g3 := NewGame("A", "B", 1)
	(DealDamage{Amount: 3, After: IfSurvives, Target: Target{Kind: TargetChosenCreature}, Then: DiscardCard{Player: ItsOwner, Zones: []Zone{Hand}, Selection: Random{}}}).Resolve(
		&EffectContext{Resolver: g3, Controller: 0},
	)
}

func TestDamageCreatureAndNeighbor(t *testing.T) {
	e := DealDamage{Spread: CreatureAndNeighbors{Amount: 3, Splash: 3, Scope: OneNeighbor}}
	if e.Text() != "choose a creature. Deal 3 damage to the chosen creature and one of its neighbors" {
		t.Errorf("text = %q", e.Text())
	}
	// A spread that cannot be misconfigured validates cleanly.
	if err := e.validate(); err != nil {
		t.Errorf("validate() = %v, want nil", err)
	}

	// The default chooser picks the left flank; its only neighbor is the middle.
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 10), 1)
	mid := g.AddToBattleline(testCreature("mid", 10), 1)
	right := g.AddToBattleline(testCreature("right", 10), 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Damage(left) != 3 || g.Damage(mid) != 3 || g.Damage(right) != 0 {
		t.Errorf("damage = %d/%d/%d, want 3/3/0", g.Damage(left), g.Damage(mid), g.Damage(right))
	}

	// A middle creature has two neighbors, so the controller picks one.
	g2 := NewGame("A", "B", 1)
	l2 := g2.AddToBattleline(testCreature("l", 10), 1)
	m2 := g2.AddToBattleline(testCreature("m", 10), 1)
	r2 := g2.AddToBattleline(testCreature("r", 10), 1)
	g2.SetChooser(0, midThenFirst{mid: m2})
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Damage(m2) != 3 || g2.Damage(l2) != 3 || g2.Damage(r2) != 0 {
		t.Errorf(
			"damage = %d/%d/%d, want mid 3 + one neighbor 3",
			g2.Damage(l2),
			g2.Damage(m2),
			g2.Damage(r2),
		)
	}

	// No creatures: nothing happens.
	g3 := NewGame("A", "B", 1)
	e.Resolve(&EffectContext{Resolver: g3, Controller: 0})
}

func TestSpreadDivideDamage(t *testing.T) {
	t.Run("text with a Per count front-loads the source", func(t *testing.T) {
		e := DealDamage{Spread: DivideDamage{
			Amount: 2,
			Per:    CardsInPlay{Player: Controller, House: namedHouse(Brobnar), Type: Creature},
		}}
		want := "deal 2 damage for each friendly Brobnar creature, " +
			"divided among any number of creatures"
		if e.Text() != want {
			t.Errorf("text = %q, want %q", e.Text(), want)
		}
	})

	t.Run("text without a Per count states the flat amount", func(t *testing.T) {
		e := DealDamage{Spread: DivideDamage{Amount: 3}}
		want := "deal 3 damage, divided among any number of creatures"
		if e.Text() != want {
			t.Errorf("text = %q, want %q", e.Text(), want)
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (DealDamage{Spread: DivideDamage{Amount: 0}}).validate() == nil {
			t.Error("a zero amount should be invalid")
		}
		if (DealDamage{Spread: DivideDamage{Amount: 1}}).validate() != nil {
			t.Error("a positive amount should be valid")
		}
	})

	t.Run("divides the pool among the chosen creatures", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 9), 0)
		b := g.AddToBattleline(testCreature("b", 9), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{a, b, a, b}})
		DealDamage{Spread: DivideDamage{Amount: 4}}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
		if g.Damage(a) != 2 || g.Damage(b) != 2 {
			t.Errorf("damage = %d/%d, want 2/2", g.Damage(a), g.Damage(b))
		}
	})

	t.Run("scales the pool by the Per count", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 9), 0)
		g.AddToBattleline(testCreature("b", 9), 0)
		// Two friendly Brobnar creatures, 2 each, all placed on the first by the
		// default chooser.
		DealDamage{Spread: DivideDamage{
			Amount: 2,
			Per:    CardsInPlay{Player: Controller, House: namedHouse(Brobnar), Type: Creature},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(a) != 4 {
			t.Errorf("damage on a = %d, want 4", g.Damage(a))
		}
	})

	t.Run("a zero pool deals nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToBattleline(testCreature("a", 9), 0)
		DealDamage{Spread: DivideDamage{
			Amount: 2,
			Per:    CardsInPlay{Player: Opponent, House: namedHouse(Brobnar), Type: Creature},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	})

	t.Run("with no creatures in play deals nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		DealDamage{Spread: DivideDamage{Amount: 3}}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	})
}
