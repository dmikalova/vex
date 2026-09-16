package engine

import "testing"

// TestFoughtCreatureIsMostPowerfulEnemy covers the Before Fight condition Baldric
// the Bold reads: unmet when there is no fought creature, unmet when an enemy
// creature has strictly greater power, and met when the fought creature ties for
// or exceeds every enemy's power (a tie still qualifies).
func TestFoughtCreatureIsMostPowerfulEnemy(t *testing.T) {
	cond := FoughtCreatureIsMostPowerfulEnemy{}
	if got := cond.CondText(); got != "if the fought creature is the most powerful enemy creature" {
		t.Errorf("CondText = %q", got)
	}

	g := started(t)
	fought := g.AddToBattleline(testCreature("fought", 4), 1)
	bigger := g.AddToBattleline(testCreature("bigger", 6), 1)

	// No fought creature in context: unmet.
	if cond.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("no fought creature: condition should be unmet")
	}

	// A stronger enemy creature exists: the fought creature is not the most powerful.
	if cond.Met(&EffectContext{Resolver: g, Controller: 0, It: fought, HasIt: true}) {
		t.Error("stronger enemy present: condition should be unmet")
	}

	// Fight the bigger one, which ties the top power: met.
	if !cond.Met(&EffectContext{Resolver: g, Controller: 0, It: bigger, HasIt: true}) {
		t.Error("most powerful enemy fought: condition should be met")
	}
}
