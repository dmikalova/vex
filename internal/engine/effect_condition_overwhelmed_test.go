package engine

import "testing"

func TestOverwhelmed(t *testing.T) {
	if (Overwhelmed{}).CondText() != "if you are overwhelmed" {
		t.Errorf("cond text = %q", (Overwhelmed{}).CondText())
	}
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 2), 1)
	g.AddToBattleline(testCreature("o2", 2), 1)
	g.AddToBattleline(testCreature("m1", 2), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !(Overwhelmed{}).Met(ctx) {
		t.Error("should be overwhelmed when the opponent controls more creatures")
	}
	g.AddToBattleline(testCreature("m2", 2), 0)
	if (Overwhelmed{}).Met(ctx) {
		t.Error("should not be overwhelmed at parity")
	}
}
