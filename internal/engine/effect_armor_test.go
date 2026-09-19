package engine

import "testing"

func TestLoseArmorValidatesItsTarget(t *testing.T) {
	if err := (LoseArmor{}).validate(); err == nil {
		t.Error("an untargeted LoseArmor should be rejected")
	}
	e := LoseArmor{Target: Target{Kind: TargetEachEnemyCreature}.WithArmor()}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if want := "each enemy creature with armor loses all of its armor"; e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
}

func TestLoseArmorStripsAndTallies(t *testing.T) {
	g := started(t)
	plated := g.AddToBattleline(
		NewCard("Plated", Dis, Creature, Common, WithPower(5), WithArmor(2)),
		1,
	)
	bare := g.AddToBattleline(testCreature("Bare", 5), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	target := Target{Kind: TargetEachEnemyCreature}.WithArmor()

	// Only the armored creature is a target at all, so only it is stripped.
	if got := target.Select(ctx); len(got) != 1 || got[0] != plated {
		t.Errorf("selected %v, want just the armored creature %d (not %d)", got, plated, bare)
	}
	LoseArmor{Target: target}.Resolve(ctx)

	if got := g.State.Cards[plated].ArmorRemaining; got != 0 {
		t.Errorf("armor remaining = %d, want 0", got)
	}
	if got := g.ArmorStripped(plated); got != 2 {
		t.Errorf("armor stripped = %d, want 2", got)
	}
	if got := ArmorLostThisWay.perTargetValue(ctx, plated); got != 2 {
		t.Errorf("per-target value = %d, want 2", got)
	}
	if want := "point of armor it lost this way"; ArmorLostThisWay.perTargetText() != want {
		t.Errorf("per-target text = %q, want %q", ArmorLostThisWay.perTargetText(), want)
	}

	// A second strip adds to the tally rather than replacing it.
	g.State.Cards[plated].ArmorRemaining = 1
	g.StripArmor(plated)
	if got := g.ArmorStripped(plated); got != 3 {
		t.Errorf("armor stripped = %d, want 3 after a second strip", got)
	}

	// The strip is turn-scoped: readying refreshes the armor and clears the tally.
	g.readyPhase(1)
	if got := g.State.Cards[plated].ArmorRemaining; got != 2 {
		t.Errorf("armor remaining = %d, want the full 2 back after readying", got)
	}
	if got := g.ArmorStripped(plated); got != 0 {
		t.Errorf("armor stripped = %d, want 0 after readying", got)
	}
}

func TestStripArmorSkipsACardOutOfPlay(t *testing.T) {
	g := started(t)
	gone := g.Register(NewCard("Gone", Dis, Creature, Common, WithPower(5), WithArmor(2)), 0)
	g.State.Discard[0].add(gone)

	g.StripArmor(gone) // must not write in-play state onto a discarded card

	if got := g.State.Cards[gone].ArmorStripped; got != 0 {
		t.Errorf("armor stripped = %d, want 0 for a card out of play", got)
	}
}

// TestAfterArmorPreventsFires exercises the After This Creature Prevents Damage
// With Its Armor trigger: a creature that spends armor absorbing combat damage
// fires its ability scaled by the amount it prevented, and one that spends no
// armor does not.
func TestAfterArmorPreventsFires(t *testing.T) {
	marucker := func() CardDefinition {
		return NewCard("marucker", Brobnar, Creature, Common,
			WithPower(5), WithArmor(2),
			WithAbility(TriggerAfterArmorPrevents, CaptureAember{
				Amount: 1,
				Per:    DamagePrevented{},
				Target: Target{Kind: TargetThisCreature},
				Source: Opponent,
			}))
	}

	t.Run("captures for each damage prevented", func(t *testing.T) {
		g := started(t)
		g.SetAember(1, 5)
		att := g.AddToBattleline(marucker(), 0)
		def := g.AddToBattleline(testCreature("def", 3), 1)
		if err := g.Fight(0, att, def); err != nil {
			t.Fatalf("Fight: %v", err)
		}
		// The 3 return damage is absorbed 2 by armor, 1 lands: 2 prevented, so 2
		// captured from the opponent's pool.
		if g.AmberOn(att) != 2 {
			t.Errorf("captured = %d, want 2", g.AmberOn(att))
		}
	})

	t.Run("does not fire when no armor is spent", func(t *testing.T) {
		g := started(t)
		g.SetAember(1, 5)
		att := g.AddToBattleline(marucker(), 0)
		def := g.AddToBattleline(testCreature("def", 0), 1)
		if err := g.Fight(0, att, def); err != nil {
			t.Fatalf("Fight: %v", err)
		}
		if g.AmberOn(att) != 0 {
			t.Errorf("captured = %d, want 0", g.AmberOn(att))
		}
	})
}

// TestDamagePreventedCount checks the DamagePrevented count reads the armor a
// creature just spent and renders its clauses.
func TestDamagePreventedCount(t *testing.T) {
	c := DamagePrevented{}
	if c.CountText() != "damage just prevented" {
		t.Errorf("CountText = %q", c.CountText())
	}
	if got := c.CountClause("3", true); got != "3 damage was just prevented" {
		t.Errorf("CountClause = %q", got)
	}
	ctx := &EffectContext{}
	if c.Value(ctx) != 0 {
		t.Errorf("Value with no prevention = %d, want 0", c.Value(ctx))
	}
	ctx.Produced.ArmorPrevented = 4
	if c.Value(ctx) != 4 {
		t.Errorf("Value = %d, want 4", c.Value(ctx))
	}
}
