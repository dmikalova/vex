package engine

import "testing"

// The tide is neutral at the start of the game, then high for the player who
// raised it and low for their opponent. No card raises it yet, so a fresh game
// reads neutral and both checks are false; the enum values still resolve to the
// right point of view once something sets them.
func TestTide(t *testing.T) {
	cases := []struct {
		tide                         Tide
		highP0, highP1, lowP0, lowP1 bool
	}{
		{TideNeutral, false, false, false, false},
		{TideHighForP0, true, false, false, true},
		{TideHighForP1, false, true, true, false},
	}
	for _, c := range cases {
		g := NewGame("A", "B", 1)
		g.State.Tide = c.tide
		if got := g.TideIsHigh(0); got != c.highP0 {
			t.Errorf("tide %d: TideIsHigh(0) = %v, want %v", c.tide, got, c.highP0)
		}
		if got := g.TideIsHigh(1); got != c.highP1 {
			t.Errorf("tide %d: TideIsHigh(1) = %v, want %v", c.tide, got, c.highP1)
		}
		if got := g.TideIsLow(0); got != c.lowP0 {
			t.Errorf("tide %d: TideIsLow(0) = %v, want %v", c.tide, got, c.lowP0)
		}
		if got := g.TideIsLow(1); got != c.lowP1 {
			t.Errorf("tide %d: TideIsLow(1) = %v, want %v", c.tide, got, c.lowP1)
		}
	}
}

// The tide conditions read from the ability's controller's point of view: with
// the tide high for player 0 it is high for controller 0 and low for controller 1.
func TestTideConditions(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Tide = TideHighForP0
	low := TideIsLow{}
	high := TideIsHigh{}

	if !high.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("TideIsHigh should be met for controller 0 while the tide is high for them")
	}
	if low.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("TideIsLow should not be met for controller 0 while the tide is high for them")
	}
	if !low.Met(&EffectContext{Resolver: g, Controller: 1}) {
		t.Error(
			"TideIsLow should be met for controller 1 while the tide is high for their opponent",
		)
	}
	if high.Met(&EffectContext{Resolver: g, Controller: 1}) {
		t.Error("TideIsHigh should not be met for controller 1 while the tide is low for them")
	}

	if got := low.CondText(); got != "if the tide is low" {
		t.Errorf("TideIsLow.CondText() = %q", got)
	}
	if got := high.CondText(); got != "if the tide is high" {
		t.Errorf("TideIsHigh.CondText() = %q", got)
	}
}
