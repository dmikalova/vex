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
	g.settleDestroyed(0) // the resolution boundary settles the buff loss (ADR 0029)

	if g.inPlay(victim) {
		t.Errorf("a creature at 0 power should have been destroyed")
	}
}

// TestArtifactSelfDestroysWhenNoCreatures covers Doom Sigil: an artifact carrying
// a DestroyedWhen condition holds while a creature is in play and destroys itself
// once the board empties. Only artifacts in play with a met condition qualify.
func TestArtifactSelfDestroysWhenNoCreatures(t *testing.T) {
	g := started(t)
	sigil := g.AddArtifact(
		NewCard("Doom Sigil", Shadows, Artifact, Rare,
			WithDestroyedWhen(InPlay{Player: EachPlayer, Type: Creature, None: true})), 0)
	creature := g.AddToBattleline(
		NewCard("Sapling", Untamed, Creature, Common, WithPower(2)), 0)

	// A creature is in play, so the artifact holds.
	g.settleDestroyed(0)
	if !g.inPlay(sigil) {
		t.Fatal("Doom Sigil should survive while a creature is in play")
	}
	// A plain artifact with no DestroyedWhen never self-destroys, and a creature is
	// never an artifact self-destroy candidate.
	plain := g.AddArtifact(NewCard("Plain", Shadows, Artifact, Common), 0)
	if g.artifactShouldSelfDestroy(plain) {
		t.Error("an artifact with no DestroyedWhen should not self-destroy")
	}
	if g.artifactShouldSelfDestroy(creature) {
		t.Error("a creature is not an artifact self-destroy candidate")
	}

	// The board empties: the artifact destroys itself, and the plain one survives.
	g.putIntoHand(creature)
	g.settleDestroyed(0)
	if g.inPlay(sigil) {
		t.Error("Doom Sigil should self-destroy once no creatures remain")
	}
	if !g.inPlay(plain) {
		t.Error("a plain artifact should survive an empty board")
	}
	// An out-of-play artifact is not a candidate.
	if g.artifactShouldSelfDestroy(sigil) {
		t.Error("a destroyed artifact should not be a self-destroy candidate")
	}
}

// TestBuffLossKillsADamagedCreature checks a damaged creature is destroyed once a
// lost buff leaves its damage at or above its remaining power.
func TestBuffLossKillsADamagedCreature(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(banner(2), 0)
	victim := g.AddToBattleline(NewCard("Oak", Untamed, Creature, Common, WithPower(3)), 0)
	g.applyRawDamage(DamageTarget{ID: victim, Amount: 4, IgnoreArmor: true})

	if !g.inPlay(victim) {
		t.Fatal("4 damage should not destroy a 5-power creature")
	}
	g.putIntoHand(src)
	g.settleDestroyed(0) // the resolution boundary settles the buff loss (ADR 0029)

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
	g.applyRawDamage(DamageTarget{ID: middle, Amount: 6, IgnoreArmor: true})
	g.putIntoHand(src)
	g.settleDestroyed(0) // the resolution boundary settles the buff loss (ADR 0029)

	if g.inPlay(middle) {
		t.Errorf("the damaged buffer should have died once its own buff was gone")
	}
	if g.inPlay(last) {
		t.Errorf("the creature the dead buffer was propping up should have died too")
	}
}
