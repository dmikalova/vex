package engine

import "testing"

// TestPlaceCounter covers Wretched Doll's marker: placing a doom counter,
// reading it back, and stacking a second onto the same entry.
func TestPlaceCounter(t *testing.T) {
	g := NewGame("A", "B", 1)
	mark := g.AddToBattleline(testCreature("mark", 3), 1)
	safe := g.AddToBattleline(testCreature("safe", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: mark, Controller: 0}

	place := PlaceCounter{Kind: CounterDoom, Target: Target{Kind: TargetEachEnemyCreature}}
	if got := (PlaceCounter{Kind: CounterDoom, Target: Target{Kind: TargetChosenCreature}}).Text(); got != "put a doom counter on a creature" {
		t.Errorf("text = %q", got)
	}
	if got := (PlaceCounter{Kind: CounterDoom, Target: Target{Kind: TargetChosenCreature}, Amount: 2}).Text(); got != "put 2 doom counters on a creature" {
		t.Errorf("plural text = %q", got)
	}
	if err := (PlaceCounter{Target: Target{Kind: TargetChosenCreature}}).validate(); err == nil {
		t.Error("PlaceCounter without a kind should not validate")
	}
	if err := (PlaceCounter{Kind: CounterDoom}).validate(); err == nil {
		t.Error("PlaceCounter without a target should not validate")
	}
	if err := place.validate(); err != nil {
		t.Errorf("validate: %v", err)
	}

	place.Resolve(ctx)
	if g.CountersOn(mark, CounterDoom) != 1 || g.CountersOn(safe, CounterDoom) != 1 {
		t.Error("both enemy creatures should carry a doom counter")
	}

	// A second placement stacks into the same entry.
	place.Resolve(ctx)
	if g.CountersOn(mark, CounterDoom) != 2 {
		t.Errorf("doom counters should stack, got %d", g.CountersOn(mark, CounterDoom))
	}
	if g.State.CounterCount != 2 {
		t.Errorf("two creatures should hold one entry each, got %d", g.State.CounterCount)
	}
}

// TestCountersOnMissing covers placing on and reading a counter off a card that
// is not in play, and a non-positive amount.
func TestCountersOnMissing(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.PlaceCounter(LocalID(0), CounterDoom, 1) // no such card: silently ignored
	if g.CountersOn(LocalID(0), CounterDoom) != 0 {
		t.Error("a card that is not in play carries no doom counter")
	}
	c := g.AddToBattleline(testCreature("c", 3), 1)
	g.PlaceCounter(c, CounterDoom, 0) // non-positive amount: no-op
	if g.CountersOn(c, CounterDoom) != 0 {
		t.Error("placing zero counters should be a no-op")
	}
}

// TestCountersShedOnLeavePlay covers that a card sheds its counters when it
// leaves play, compacting the global table.
func TestCountersShedOnLeavePlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	victim := g.AddToBattleline(testCreature("victim", 3), 1)
	other := g.AddToBattleline(testCreature("other", 3), 1)
	g.PlaceCounter(victim, CounterDoom, 1)
	g.PlaceCounter(other, CounterDoom, 1)

	ctx := &EffectContext{Resolver: g, Source: victim, Controller: 0}
	Destroy{Target: Target{Kind: TargetChosenCreature}.WithCounter(CounterDoom)}.Resolve(ctx)

	if g.State.CounterCount != 1 {
		t.Fatalf(
			"the destroyed creature's entry should be shed, got %d entries",
			g.State.CounterCount,
		)
	}
	if g.CountersOn(other, CounterDoom) != 1 {
		t.Error("the surviving creature should keep its counter")
	}
}

// TestCounterTableOverflow covers the caught invariant when a 65th distinct
// (card, kind) pair is placed.
func TestCounterTableOverflow(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("overflow", 3), 1)
	// Fill the table with entries for cards that cannot collide with c's LocalID.
	for i := 0; i < maxCounterEntries; i++ {
		g.State.Counters[i] = CounterEntry{Card: LocalID(128 + i), Kind: CounterDoom, N: 1}
	}
	g.State.CounterCount = maxCounterEntries

	defer func() {
		if recover() == nil {
			t.Error("placing a 65th distinct counter pair should panic")
		}
	}()
	g.PlaceCounter(c, CounterDoom, 1)
}

// TestCounterInPlay covers the condition Wretched Doll checks before it destroys
// or marks.
func TestCounterInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: c, Controller: 0}

	cond := CounterInPlay{Kind: CounterDoom}
	if cond.CondText() != "if there is a doom counter in play" {
		t.Errorf("cond text = %q", cond.CondText())
	}
	if cond.Met(ctx) {
		t.Error("no doom counter in play yet")
	}

	g.PlaceCounter(c, CounterDoom, 1)
	if !cond.Met(ctx) {
		t.Error("a doom counter is now in play")
	}
}

// TestDestroyWithCounter covers the Target filter: destroying only the creatures
// that carry a doom counter.
func TestDestroyWithCounter(t *testing.T) {
	g := NewGame("A", "B", 1)
	doomed := g.AddToBattleline(testCreature("doomed", 3), 1)
	safe := g.AddToBattleline(testCreature("safe", 3), 1)
	g.PlaceCounter(doomed, CounterDoom, 1)
	ctx := &EffectContext{Resolver: g, Source: doomed, Controller: 0}

	target := Target{Kind: TargetEachCreature}.WithCounter(CounterDoom)
	if got := target.Text(); got != "each creature with a doom counter" {
		t.Errorf("target text = %q", got)
	}

	Destroy{Target: target}.Resolve(ctx)
	if g.InPlay(doomed) {
		t.Error("the doomed creature should be destroyed")
	}
	if !g.InPlay(safe) {
		t.Error("the unmarked creature should survive")
	}
}

// TestCounterKindNoun covers the fallback noun and validity of the zero kind.
func TestCounterKindNoun(t *testing.T) {
	if got := CounterNone.noun(); got != "counter" {
		t.Errorf("fallback noun = %q", got)
	}
	if CounterNone.valid() {
		t.Error("CounterNone should not be a valid kind")
	}
}

// TestCounterSaturates covers the uint8 clamp when a huge count is placed.
func TestCounterSaturates(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 1)
	g.PlaceCounter(c, CounterDoom, 300)
	if g.CountersOn(c, CounterDoom) != 255 {
		t.Errorf("counter should saturate at 255, got %d", g.CountersOn(c, CounterDoom))
	}
	g.PlaceCounter(c, CounterDoom, 10)
	if g.CountersOn(c, CounterDoom) != 255 {
		t.Errorf("saturated counter should stay at 255, got %d", g.CountersOn(c, CounterDoom))
	}
}
