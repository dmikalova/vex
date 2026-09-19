package engine

import (
	"errors"
	"slices"
	"testing"
)

// TestRuleOfSixSharedByCardName confirms the six-usage pool is keyed by card name
// alone: every copy of one name draws from the same pool, including an opponent's
// copy the active player controls, since only the active player uses cards.
func TestRuleOfSixSharedByCardName(t *testing.T) {
	g := started(t)
	a := g.AddToBattleline(testCreature("Twin", 3), 0)
	b := g.AddToBattleline(testCreature("Twin", 3), 0)
	enemy := g.AddToBattleline(testCreature("Twin", 3), 1)

	g.recordUsage(a)
	g.recordUsage(a)
	g.recordUsage(b)
	if got := g.nameUsagesThisTurn(a); got != 3 {
		t.Errorf("shared name pool = %d, want 3", got)
	}
	if got := g.nameUsagesThisTurn(enemy); got != 3 {
		t.Errorf("opponent's copy shares the same pool = %d, want 3", got)
	}
	if g.atRuleOfSix(a) {
		t.Error("three usages should not reach the Rule of Six")
	}
	g.recordUsage(a)
	g.recordUsage(a)
	g.recordUsage(enemy)
	if !g.atRuleOfSix(b) {
		t.Error("six usages of a name should reach the Rule of Six")
	}
}

// TestRuleOfSixBarsPlay confirms a player cannot play a card whose name has
// already been used six times this turn.
func TestRuleOfSixBarsPlay(t *testing.T) {
	g := started(t)
	id := g.AddToHand(testCreature("Grunt", 3), 0)
	for range RuleOfSix {
		g.recordUsage(id)
	}
	if err := g.CanPlay(0, id); !errors.Is(err, ErrRuleOfSix) {
		t.Errorf("CanPlay at the Rule of Six = %v, want ErrRuleOfSix", err)
	}
	if _, err := g.PlayCreature(0, 0, false); !errors.Is(err, ErrRuleOfSix) {
		t.Errorf("PlayCreature at the Rule of Six = %v, want ErrRuleOfSix", err)
	}
}

// TestRuleOfSixBarsUse confirms a creature cannot be used (reaped, fought, or its
// Action: fired) once its name has been used six times this turn.
func TestRuleOfSixBarsUse(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("Worker", 3), 0)
	g.State.Cards[id].Exhausted = false
	if err := g.usable(0, id); err != nil {
		t.Fatalf("a fresh creature should be usable, got %v", err)
	}
	for range RuleOfSix {
		g.recordUsage(id)
	}
	if err := g.usable(0, id); !errors.Is(err, ErrRuleOfSix) {
		t.Errorf("use at the Rule of Six = %v, want ErrRuleOfSix", err)
	}
}

// TestRuleOfSixBarsDestroyedResolution confirms a Destroyed: ability does not
// resolve once its creature's name has been used six times this turn, though the
// creature still leaves play.
func TestRuleOfSixBarsDestroyedResolution(t *testing.T) {
	g := started(t)
	dying := testCreature("Martyr", 3, WithAbility(
		TriggerDestroyed, GainAember{Amount: 1, Player: Controller}))
	id := g.AddToBattleline(dying, 0)
	for range RuleOfSix {
		g.recordUsage(id)
	}

	g.DestroyEach(0, []LocalID{id})

	if g.Aember(0) != 0 {
		t.Errorf("a barred Destroyed: ability gained %d Æmber, want 0", g.Aember(0))
	}
	if slices.Contains(g.Battleline(0), id) {
		t.Error("the creature should still leave play after its Destroyed: is barred")
	}
}

// TestRuleOfSixBarsAbilityDrivenUse confirms an ability that uses a creature
// (Legatus Raptor's "ready and use another friendly creature") does nothing once
// that creature's name has been used six times this turn — so two Legatus Raptors
// sharing the "Raptor" pool cannot ready-and-use each other past the six between
// them.
func TestRuleOfSixBarsAbilityDrivenUse(t *testing.T) {
	g := started(t)
	reaper := g.AddToBattleline(testCreature("Raptor", 3), 0)
	fighter := g.AddToBattleline(testCreature("Raptor", 3), 0)
	enemy := g.AddToBattleline(testCreature("Prey", 3), 1)
	for range RuleOfSix {
		g.recordUsage(reaper) // fills the shared "Raptor" pool
	}

	g.ReapWith(reaper)
	if g.Aember(0) != 0 {
		t.Errorf("an ability-driven reap at the Rule of Six gained %d Æmber, want 0", g.Aember(0))
	}
	if g.Exhausted(reaper) {
		t.Error("a creature at the Rule of Six should not exhaust from an ability-driven reap")
	}

	g.FightWith(fighter, enemy)
	if g.Exhausted(fighter) {
		t.Error("a creature at the Rule of Six should not fight from an ability")
	}
	if g.Damage(enemy) != 0 {
		t.Errorf("a fight barred by the Rule of Six dealt %d damage, want 0", g.Damage(enemy))
	}
}
