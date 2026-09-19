package engine

import "testing"

func TestHaunted(t *testing.T) {
	c := Haunted{}
	if got := c.CondText(); got != "if you are haunted" {
		t.Errorf("CondText() = %q, want %q", got, "if you are haunted")
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if c.Met(ctx) {
		t.Error("empty discard pile should not be haunted")
	}

	for range 9 {
		g.AddToDiscard(testCreature("buried", 1), 0)
	}
	if c.Met(ctx) {
		t.Error("9 cards in discard should not be haunted")
	}

	g.AddToDiscard(testCreature("buried", 1), 0)
	if !c.Met(ctx) {
		t.Error("10 cards in discard should be haunted")
	}

	// The opponent's discard pile does not haunt the controller.
	g2 := NewGame("A", "B", 1)
	for range 10 {
		g2.AddToDiscard(testCreature("theirs", 1), 1)
	}
	if c.Met(&EffectContext{Resolver: g2, Controller: 0}) {
		t.Error("opponent's discard pile should not haunt the controller")
	}
}
