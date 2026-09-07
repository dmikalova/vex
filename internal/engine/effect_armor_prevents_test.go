package engine

import "testing"

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
