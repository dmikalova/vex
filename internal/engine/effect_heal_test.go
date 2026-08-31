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

func TestDealDamageAmountFrom(t *testing.T) {
	// Text and validate for the "deal that amount of damage" mode.
	e := DealDamage{AmountFrom: DamageHealed{}, Target: Target{Kind: TargetChosenOtherCreature}}
	if got := e.Text(); got != "deal that amount of damage to another creature" {
		t.Errorf("text = %q", got)
	}
	if (DealDamage{Target: Target{Kind: TargetEachCreature}, AmountFrom: DamageHealed{}, Per: DamageHealed{}}).validate() == nil {
		t.Error("AmountFrom + Per together should be invalid")
	}
	if (DamageHealed{}).CountText() != "damage healed this way" {
		t.Errorf("count text = %q", (DamageHealed{}).CountText())
	}

	// Guardian Demon's composition: heal a creature (records the amount), then deal
	// that much to another creature. A fixed Heal of 2 caps at the creature's damage.
	g := NewGame("A", "B", 1)
	wounded := g.AddToBattleline(testCreature("wounded", 5), 0)
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	g.State.Cards[wounded].Damage = 3
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{wounded, foe}})
	ctx := &EffectContext{Resolver: g, Controller: 0}
	Heal{Amount: 2, Target: Target{Kind: TargetChosenCreature}}.Resolve(ctx)
	if g.Damage(wounded) != 1 || ctx.Produced.DamageHealed != 2 {
		t.Errorf(
			"after heal: damage=%d healed=%d, want 1/2",
			g.Damage(wounded),
			ctx.Produced.DamageHealed,
		)
	}
	e.Resolve(ctx) // ctx.It = wounded, so OtherCreature offers foe
	if g.Damage(foe) != 2 {
		t.Errorf("foe damage = %d, want 2 (that amount)", g.Damage(foe))
	}

	// A fixed Heal heals only as much damage as the creature has (clamped).
	g2 := NewGame("A", "B", 1)
	c := g2.AddToBattleline(testCreature("c", 5), 0)
	g2.State.Cards[c].Damage = 1
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	Heal{Amount: 2, Target: Target{Kind: TargetEachFriendlyCreature}}.Resolve(ctx2)
	if g2.Damage(c) != 0 || ctx2.Produced.DamageHealed != 1 {
		t.Errorf(
			"clamped heal: damage=%d healed=%d, want 0/1",
			g2.Damage(c),
			ctx2.Produced.DamageHealed,
		)
	}

	// Without a context creature, OtherCreature offers every creature.
	g3 := NewGame("A", "B", 1)
	a := g3.AddToBattleline(testCreature("a", 3), 0)
	g3.AddToBattleline(testCreature("b", 3), 1)
	g3.SetChooser(0, &idQueueChooser{ids: []LocalID{a}})
	if ids := (Target{Kind: TargetChosenOtherCreature}).Select(
		&EffectContext{Resolver: g3, Controller: 0},
	); len(ids) != 1 ||
		ids[0] != a {
		t.Errorf("no-context other-creature select = %v, want [%d]", ids, a)
	}
}

func TestSetDamageClampsNegative(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 5), 0)
	g.SetDamage(c, -3)
	if g.Damage(c) != 0 {
		t.Errorf("damage = %d, want 0 (clamped)", g.Damage(c))
	}
}
