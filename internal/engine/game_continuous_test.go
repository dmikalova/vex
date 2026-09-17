package engine

import "testing"

// TestSetDamageImmuneSingleCreature covers a per-creature immunity: only the named
// creature is protected, not its neighbor (the ScopeSubject reach).
func TestSetDamageImmuneSingleCreature(t *testing.T) {
	g := started(t)
	guarded := g.AddToBattleline(testCreature("guarded", 5), 0)
	other := g.AddToBattleline(testCreature("other", 5), 0)

	g.SetDamageImmune(guarded, RemainderOfPlayerTurn)

	if !g.DamageImmune(guarded) {
		t.Error("the named creature should be immune")
	}
	if g.DamageImmune(other) {
		t.Error("a different creature should not be immune")
	}
}

// TestAddContinuousOverflowPanics checks the stack panics rather than silently
// dropping an effect past capacity.
func TestAddContinuousOverflowPanics(t *testing.T) {
	g := started(t)
	defer func() {
		if recover() == nil {
			t.Error("filling the continuous stack past capacity should panic")
		}
	}()
	for i := 0; i <= maxContinuous; i++ {
		g.addContinuous(ContinuousEffect{
			Kind:  ContinuousDamageImmune,
			Scope: ScopeAllCreatures,
		}, RemainderOfPlayerTurn)
	}
}
