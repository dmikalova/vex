package engine

import "testing"

// TestSourceIsFightingCondition covers the "while fighting" condition and the
// CurrentlyFighting flag it reads: unmet outside combat, met while a creature is
// one of the fight's combatants.
func TestSourceIsFightingCondition(t *testing.T) {
	if got := (SourceIsFighting{}).CondText(); got != "if fighting" {
		t.Errorf("CondText = %q, want %q", got, "if fighting")
	}
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	other := g.AddToBattleline(testCreature("other", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: c}
	if (SourceIsFighting{}).Met(ctx) {
		t.Error("not fighting: condition should be unmet")
	}
	if g.CurrentlyFighting(c) {
		t.Error("CurrentlyFighting should be false with no fight in progress")
	}
	// Mark the pair as fighting, matching on the second slot to exercise it too.
	g.State.FightersPlus = [2]LocalID{other + 1, c + 1}
	if !(SourceIsFighting{}).Met(ctx) {
		t.Error("fighting: condition should be met")
	}
	if !g.CurrentlyFighting(c) {
		t.Error("CurrentlyFighting should report the marked fighter")
	}
}

// nizakLike builds a creature with Nizak, The Forgotten's two abilities: a
// combat-scoped invulnerability self-grant and a reaction that returns an enemy
// destroyed fighting it to its owner's hand.
func nizakLike(power int) CardDefinition {
	return testCreature("nizak", power,
		WithConstantAbility(ConstantAbility{
			Target:         Target{Kind: TargetThisCreature},
			Keywords:       []Keyword{Invulnerable},
			WhileCondition: SourceIsFighting{},
		}),
		WithAbility(TriggerAfterDestroyedFighting, ReturnItToHand{}),
	)
}

// TestInvulnerableWhileFighting checks the combat-scoped grant: while fighting the
// creature takes no damage and survives a fight that would otherwise destroy it,
// and the enemy it kills returns to its owner's hand; outside combat the grant is
// gone.
func TestInvulnerableWhileFighting(t *testing.T) {
	g := started(t)
	nizak := g.AddToBattleline(nizakLike(6), 0)
	prey := g.AddToBattleline(testCreature("prey", 6), 1)

	// Not fighting yet: the self-grant is suspended.
	if g.hasKeyword(nizak, Invulnerable) {
		t.Fatal("should not be invulnerable before the fight")
	}

	if err := g.Fight(0, nizak, prey); err != nil {
		t.Fatalf("Fight: %v", err)
	}

	// Equal power would trade both creatures, but invulnerability spares the
	// attacker: it survives undamaged.
	if !g.inPlay(nizak) {
		t.Fatal("invulnerable-while-fighting creature should survive the fight")
	}
	if got := g.State.Cards[nizak].Damage; got != 0 {
		t.Errorf("damage while fighting = %d, want 0", got)
	}
	// The destroyed enemy returned to its owner's hand rather than the discard.
	if g.inPlay(prey) {
		t.Fatal("prey should have left play")
	}
	inHand := false
	for _, id := range g.Hand(1) {
		if id == prey {
			inHand = true
		}
	}
	if !inHand {
		t.Error("prey should be in its owner's hand")
	}
	// The fight is over: the grant lifts and it is no longer invulnerable.
	if g.CurrentlyFighting(nizak) {
		t.Error("should not be fighting after the fight resolves")
	}
	if g.hasKeyword(nizak, Invulnerable) {
		t.Error("should not be invulnerable outside combat")
	}
}
