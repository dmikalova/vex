package engine

import "testing"

// TestAfterAemberStolenFromYouFires exercises the After Æmber Is Stolen From You
// trigger: a theft from the victim fires their in-play ability scaled by the
// amount stolen in that theft, while a theft from the other player does not.
func TestAfterAemberStolenFromYouFires(t *testing.T) {
	molephin := func() CardDefinition {
		return NewCard("molephin", Untamed, Creature, Common,
			WithPower(3),
			WithAbility(TriggerAfterAemberStolenFromYou, DealDamage{
				Amount: 1,
				Per:    AemberStolenThisEvent{},
				Target: Target{Kind: TargetEachEnemyCreature},
			}))
	}

	t.Run("fires scaled by the amount stolen from you", func(t *testing.T) {
		g := started(t)
		g.SetAember(1, 5)
		g.AddToBattleline(molephin(), 1)
		foe := g.AddToBattleline(testCreature("foe", 6), 0)
		// Player 0 (thief) steals 2 Æmber from player 1, the victim who controls
		// Molephin; each of player 1's enemy creatures (player 0's) takes 2 damage.
		StealAember{Amount: 2}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(foe) != 2 {
			t.Errorf("damage to enemy creature = %d, want 2", g.Damage(foe))
		}
	})

	t.Run("does not fire when the other player is robbed", func(t *testing.T) {
		g := started(t)
		g.SetAember(0, 5)
		g.AddToBattleline(molephin(), 1)
		foe := g.AddToBattleline(testCreature("foe", 6), 0)
		// Player 1 steals from player 0, so nothing is stolen from Molephin's
		// controller and its ability does not fire.
		StealAember{Amount: 2}.Resolve(&EffectContext{Resolver: g, Controller: 1})
		if g.Damage(foe) != 0 {
			t.Errorf("damage to enemy creature = %d, want 0", g.Damage(foe))
		}
	})

	t.Run("a theft of nothing fires nothing", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(molephin(), 1)
		foe := g.AddToBattleline(testCreature("foe", 6), 0)
		g.EmitAemberStolenFrom(1, 0)
		if g.Damage(foe) != 0 {
			t.Errorf("damage to enemy creature = %d, want 0", g.Damage(foe))
		}
	})
}

// TestAemberStolenThisEventCount checks the AemberStolenThisEvent count reads the
// amount stashed on the context and renders its clauses.
func TestAemberStolenThisEventCount(t *testing.T) {
	c := AemberStolenThisEvent{}
	if c.CountText() != "Æmber stolen" {
		t.Errorf("CountText = %q", c.CountText())
	}
	if got := c.CountClause("2", true); got != "2 Æmber was stolen" {
		t.Errorf("CountClause = %q", got)
	}
	ctx := &EffectContext{}
	if c.Value(ctx) != 0 {
		t.Errorf("Value with no theft = %d, want 0", c.Value(ctx))
	}
	ctx.Produced.AemberStolen = 3
	if c.Value(ctx) != 3 {
		t.Errorf("Value = %d, want 3", c.Value(ctx))
	}
}
