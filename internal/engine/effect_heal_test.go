package engine

import "testing"

func TestHealEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 5), 0)
	other := g.AddToBattleline(testCreature("other", 5), 0)
	g.State.Cards[src].Damage = 3
	g.State.Cards[other].Damage = 4
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	partial := Heal{Amount: 2, Target: Target{Kind: TargetThisCreature}}
	if partial.Text() != "heal 2 damage from "+SelfName {
		t.Errorf("partial text = %q", partial.Text())
	}
	partial.Resolve(ctx)
	if g.State.Cards[src].Damage != 1 {
		t.Errorf("partial heal: src damage = %d, want 1", g.State.Cards[src].Damage)
	}
	(Heal{Amount: 5, Target: Target{Kind: TargetThisCreature}}).Resolve(
		ctx,
	) // over-heal floors at 0
	if g.State.Cards[src].Damage != 0 {
		t.Errorf("over-heal should floor at 0, got %d", g.State.Cards[src].Damage)
	}

	full := Heal{Fully: true, Target: Target{Kind: TargetEachOtherFriendlyCreature}}
	if full.Text() != "fully heal each other friendly creature" {
		t.Errorf("full text = %q", full.Text())
	}
	full.Resolve(ctx)
	if g.State.Cards[other].Damage != 0 {
		t.Errorf("full heal: other damage = %d, want 0", g.State.Cards[other].Damage)
	}
}

func TestHealValidate(t *testing.T) {
	this := Target{Kind: TargetThisCreature}
	if err := (Heal{Target: this, Fully: true, Amount: 2}).validate(); err == nil {
		t.Error("Heal with both Amount and Fully should be invalid")
	}
	if err := (Heal{Target: this, Fully: true}).validate(); err != nil {
		t.Errorf("full heal should be valid, got %v", err)
	}
	if err := (Heal{Target: this, Amount: 2}).validate(); err != nil {
		t.Errorf("fixed heal should be valid, got %v", err)
	}
	if err := (Heal{Fully: true}).validate(); err == nil {
		t.Error("an unset target should be rejected")
	}
}

func TestHealCountsCreaturesHealed(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 5), 0)
	b := g.AddToBattleline(testCreature("b", 5), 1)
	healthy := g.AddToBattleline(testCreature("healthy", 5), 0)
	g.State.Cards[a].Damage = 2
	g.State.Cards[b].Damage = 1
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Heal each creature: the two damaged ones are healed, the undamaged one is
	// skipped, and the tally lands on the context.
	Heal{Amount: 1, Target: Target{Kind: TargetEachCreature}}.Resolve(ctx)
	if g.Damage(a) != 1 || g.Damage(b) != 0 {
		t.Errorf("damage a=%d b=%d, want 1/0", g.Damage(a), g.Damage(b))
	}
	if g.Damage(healthy) != 0 {
		t.Errorf("undamaged creature should stay at 0, got %d", g.Damage(healthy))
	}
	if ctx.Produced.Healed != 2 {
		t.Errorf("ctx.Produced.Healed = %d, want 2", ctx.Produced.Healed)
	}

	// CreaturesHealed reads that tally and renders the "for each" clause.
	cnt := CreaturesHealed{}
	if cnt.Value(ctx) != 2 {
		t.Errorf("CreaturesHealed.Value = %d, want 2", cnt.Value(ctx))
	}
	if cnt.CountText() != "creature healed this way" {
		t.Errorf("CountText = %q", cnt.CountText())
	}

	// No damaged creatures: nothing is healed and the tally is zero.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(testCreature("c", 5), 0)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	Heal{Amount: 1, Target: Target{Kind: TargetEachCreature}}.Resolve(ctx2)
	if ctx2.Produced.Healed != 0 {
		t.Errorf("ctx.Produced.Healed = %d, want 0 (nothing to heal)", ctx2.Produced.Healed)
	}
}

// TestHealedContextIsolatedAcrossNestedAbilities proves each ability resolution
// gets its own EffectContext. When one ability's effect chains into a second
// ability, the nested Heal tallies onto the nested context, so the outer
// CreaturesHealed keeps reading the outer Heal's total. triggerAbilities builds a
// fresh &EffectContext per ability, so a chained card can never overwrite the
// count a parent ability is still counting against.
func TestHealedContextIsolatedAcrossNestedAbilities(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 5), 0)
	b := g.AddToBattleline(testCreature("b", 5), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 5), 1)
	g.State.Cards[a].Damage = 2
	g.State.Cards[b].Damage = 2
	g.State.Cards[enemy].Damage = 2

	// A separate friendly creature whose Play ability heals the enemy side. Firing
	// it mid-sequence heals one creature, so its own context tallies 1 — a value
	// that must not leak into the outer ability's count of 2.
	nested := g.AddToBattleline(testCreature(
		"nested",
		5,
		WithAbility(
			TriggerAfterPlay,
			Heal{Amount: 1, Target: Target{Kind: TargetEachEnemyCreature}},
		),
	), 0)

	ctx := &EffectContext{Resolver: g, Controller: 0}
	seq := Sequence{Effects: []Effect{
		Heal{
			Amount: 1,
			Target: Target{Kind: TargetEachFriendlyCreature},
		}, // heals a, b -> ctx.Produced.Healed = 2
		gameEffect{
			fn: func() { g.triggerAbilities(nested, TriggerAfterPlay, 0, false) },
		}, // nested heals enemy on its OWN context (Healed = 1)
		GainAember{
			Player: Controller,
			Amount: 1,
			Per:    CreaturesHealed{},
		}, // reads THIS context's 2, not the nested 1
	}}
	seq.Resolve(ctx)

	if ctx.Produced.Healed != 2 {
		t.Errorf(
			"outer ctx.Produced.Healed = %d, want 2 (a nested ability must not overwrite it)",
			ctx.Produced.Healed,
		)
	}
	if g.Aember(0) != 2 {
		t.Errorf("gained %d Æmber, want 2 (outer Heal's count, not the nested 1)", g.Aember(0))
	}
}

func TestHealThenDamage(t *testing.T) {
	t.Run("heals a creature then deals the healed amount to another", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		wounded := g.AddToBattleline(testCreature("wounded", 5), 0)
		foe := g.AddToBattleline(testCreature("foe", 5), 1)
		g.State.Cards[wounded].Damage = 3
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{wounded, foe}})
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := HealThenDamage{Amount: 2}
		if e.Text() != "heal up to 2 damage from a creature. Deal that amount of damage to another creature" {
			t.Errorf("text = %q", e.Text())
		}
		e.Resolve(ctx)
		if g.Damage(wounded) != 1 {
			t.Errorf("wounded damage = %d, want 1 (healed 2)", g.Damage(wounded))
		}
		if g.Damage(foe) != 2 {
			t.Errorf("foe damage = %d, want 2", g.Damage(foe))
		}
	})

	t.Run("heals and deals only as much as the creature's damage", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		wounded := g.AddToBattleline(testCreature("wounded", 5), 0)
		foe := g.AddToBattleline(testCreature("foe", 5), 1)
		g.State.Cards[wounded].Damage = 1
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{wounded, foe}})
		HealThenDamage{Amount: 2}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(wounded) != 0 || g.Damage(foe) != 1 {
			t.Errorf("damage = %d/%d, want 0/1", g.Damage(wounded), g.Damage(foe))
		}
	})

	t.Run("healing nothing deals nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		healthy := g.AddToBattleline(testCreature("healthy", 5), 0)
		foe := g.AddToBattleline(testCreature("foe", 5), 1)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{healthy}})
		HealThenDamage{Amount: 2}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(foe) != 0 {
			t.Errorf("foe damage = %d, want 0 (nothing healed)", g.Damage(foe))
		}
	})

	t.Run("with no other creature to damage, only heals", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		wounded := g.AddToBattleline(testCreature("wounded", 5), 0)
		g.State.Cards[wounded].Damage = 2
		HealThenDamage{Amount: 2}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(wounded) != 0 {
			t.Errorf("wounded damage = %d, want 0 (healed)", g.Damage(wounded))
		}
	})

	t.Run("with no creatures, does nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		HealThenDamage{Amount: 2}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	})
}
