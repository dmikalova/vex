package engine

import "testing"

// banner returns a creature whose constant ability gives every friendly creature
// the given power bonus.
func banner(bonus int) CardDefinition {
	return NewCard("Banner", Untamed, Creature, Common, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetEachFriendlyCreature},
			PowerBonus: bonus,
		}))
}

// TestZeroPowerIsDestroyed checks a creature left at 0 power by the loss of a
// buff is destroyed, without anything dealing it damage.
func TestZeroPowerIsDestroyed(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(banner(2), 0)
	victim := g.AddToBattleline(NewCard("Sapling", Untamed, Creature, Common), 0)

	if got := g.Power(victim); got != 2 {
		t.Fatalf("power = %d, want 2 while the banner is in play", got)
	}
	g.putIntoHand(src)

	if g.inPlay(victim) {
		t.Errorf("a creature at 0 power should have been destroyed")
	}
}

// TestBuffLossKillsADamagedCreature checks a damaged creature is destroyed once a
// lost buff leaves its damage at or above its remaining power.
func TestBuffLossKillsADamagedCreature(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(banner(2), 0)
	victim := g.AddToBattleline(NewCard("Oak", Untamed, Creature, Common, WithPower(3)), 0)
	g.applyRawDamage(victim, 4, true)

	if !g.inPlay(victim) {
		t.Fatal("4 damage should not destroy a 5-power creature")
	}
	g.putIntoHand(src)

	if g.inPlay(victim) {
		t.Errorf("damage at or above the remaining power should destroy the creature")
	}
}

// TestSettleCascades checks the sweep repeats: the creature it destroys was
// itself buffing another, which then dies in the same settling.
func TestSettleCascades(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(banner(2), 0)
	middle := g.AddToBattleline(banner(2), 0)
	last := g.AddToBattleline(NewCard("Sprout", Untamed, Creature, Common), 0)

	// middle sits at 3 printed + 2 from src + 2 from itself; last is 0 + 4.
	g.applyRawDamage(middle, 6, true)
	g.putIntoHand(src)

	if g.inPlay(middle) {
		t.Errorf("the damaged buffer should have died once its own buff was gone")
	}
	if g.inPlay(last) {
		t.Errorf("the creature the dead buffer was propping up should have died too")
	}
}
