package engine

import "testing"

// TestCountersOnAttachedUpgrade covers Disruption Field: a counter placed on an
// attached upgrade is stored, raises the opponent's key cost through the
// upgrade's own WithKeyCost scaled by CountersOnThis, and is shed when the
// upgrade leaves play with its host.
func TestCountersOnAttachedUpgrade(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	up := attachUpgrade(g, host, NewCard("Field", StarAlliance, Upgrade, Rare,
		WithKeyCost(NewKeyCostChange(Opponent, 1).Per(CountersOnThis{Kind: CounterDisruption}))))

	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Fatalf("opponent key cost with no counters = %d, want %d", got, KeyCost)
	}

	g.PlaceCounter(up, CounterDisruption, 2)
	if got := g.CountersOn(up, CounterDisruption); got != 2 {
		t.Fatalf("counters on the attached upgrade = %d, want 2", got)
	}
	if got := g.CurrentKeyCost(1); got != KeyCost+2 {
		t.Errorf("opponent key cost with two counters = %d, want %d", got, KeyCost+2)
	}

	ctx := &EffectContext{Resolver: g, Source: host, Controller: 0}
	Destroy{Target: Target{Kind: TargetThisCreature}}.Resolve(ctx)
	if g.State.CounterCount != 0 {
		t.Errorf(
			"the discarded upgrade should shed its counters, got %d entries",
			g.State.CounterCount,
		)
	}
}

// TestPlaceCounter covers Wretched Doll's marker: placing a doom counter,
// reading it back, and stacking a second onto the same entry.
func TestPlaceCounter(t *testing.T) {
	g := NewGame("A", "B", 1)
	mark := g.AddToBattleline(testCreature("mark", 3), 1)
	safe := g.AddToBattleline(testCreature("safe", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: mark, Controller: 0}

	place := PlaceCounter{
		Amount: 1,
		Kind:   CounterDoom,
		Target: Target{Kind: TargetEachEnemyCreature},
	}
	if got := (PlaceCounter{Amount: 1, Kind: CounterDoom, Target: Target{Kind: TargetChosenCreature}}).Text(); got != "put a doom counter on a creature" {
		t.Errorf("text = %q", got)
	}
	if got := (PlaceCounter{Kind: CounterDoom, Target: Target{Kind: TargetChosenCreature}, Amount: 2}).Text(); got != "put 2 doom counters on a creature" {
		t.Errorf("plural text = %q", got)
	}
	if err := (PlaceCounter{Amount: 1, Target: Target{Kind: TargetChosenCreature}}).validate(); err == nil {
		t.Error("PlaceCounter without a kind should not validate")
	}
	if err := (PlaceCounter{Amount: 1, Kind: CounterDoom}).validate(); err == nil {
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

// TestCountersOnThisAtLeast covers the threshold condition The Big One checks: it
// is met once the source card carries at least N counters of the kind.
func TestCountersOnThisAtLeast(t *testing.T) {
	g := NewGame("A", "B", 1)
	bomb := g.AddArtifact(NewCard("The Big One", Brobnar, Artifact, Rare), 0)
	ctx := &EffectContext{Resolver: g, Source: bomb, Controller: 0}

	cond := CountersOnThisAtLeast{Kind: CounterFuse, N: 10}
	if want := "if there are 10 or more fuse counters on " + SelfName; cond.CondText() != want {
		t.Errorf("cond text = %q, want %q", cond.CondText(), want)
	}
	if cond.Met(ctx) {
		t.Error("no fuse counters yet, should not be met")
	}

	g.PlaceCounter(bomb, CounterFuse, 9)
	if cond.Met(ctx) {
		t.Error("nine fuse counters is below the threshold of ten")
	}

	g.PlaceCounter(bomb, CounterFuse, 1)
	if !cond.Met(ctx) {
		t.Error("ten fuse counters should meet the threshold")
	}
}

// TestCounterFuseNoun covers the fuse counter's rendered noun.
func TestCounterFuseNoun(t *testing.T) {
	if got := CounterFuse.noun(); got != "fuse counter" {
		t.Errorf("CounterFuse.noun() = %q, want %q", got, "fuse counter")
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
	if got := CounterGlory.noun(); got != "glory counter" {
		t.Errorf("CounterGlory.noun() = %q, want %q", got, "glory counter")
	}
	if got := CounterDisruption.noun(); got != "disruption counter" {
		t.Errorf("CounterDisruption.noun() = %q, want %q", got, "disruption counter")
	}
	if got := CounterScheme.noun(); got != "scheme counter" {
		t.Errorf("CounterScheme.noun() = %q, want %q", got, "scheme counter")
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

// TestCountersOnThis covers the Count that reads growth counters off the source
// card: its value tracks the counters placed, and its "for each" noun names the
// counter.
func TestCountersOnThis(t *testing.T) {
	g := NewGame("A", "B", 1)
	tree := g.AddArtifact(NewCard("Vineapple Tree", Untamed, Artifact, Rare), 0)
	c := CountersOnThis{Kind: CounterGrowth}
	ctx := &EffectContext{Resolver: g, Source: tree}

	if got := c.Value(ctx); got != 0 {
		t.Errorf("value with no counters = %d, want 0", got)
	}
	g.PlaceCounter(tree, CounterGrowth, 3)
	if got := c.Value(ctx); got != 3 {
		t.Errorf("value with three counters = %d, want 3", got)
	}
	if got := c.CountText(); got != "growth counter on "+SelfName {
		t.Errorf("count text = %q", got)
	}
}

// TestGrowthCounterKeyCostText covers the rendered key-cost rule scaled by the
// source's growth counters.
func TestGrowthCounterKeyCostText(t *testing.T) {
	kc := NewKeyCostChange(EachPlayer, 1).Per(CountersOnThis{Kind: CounterGrowth})
	want := "Each player's keys cost +1 Æmber for each growth counter on " + SelfName + "."
	if got := keyCostText(kc); got != want {
		t.Errorf("key cost text = %q, want %q", got, want)
	}
}

// TestRemoveCountersEffect covers the RemoveCounters effect: its validation, its
// text, and that it drops the named kind while leaving another kind in place.
func TestRemoveCountersEffect(t *testing.T) {
	if err := (RemoveCounters{Target: Target{Kind: TargetThisCreature}}).validate(); err == nil {
		t.Error("RemoveCounters without a kind should not validate")
	}
	if err := (RemoveCounters{Kind: CounterGrowth}).validate(); err == nil {
		t.Error("RemoveCounters without a target should not validate")
	}
	e := RemoveCounters{Kind: CounterGrowth, Target: Target{Kind: TargetThisCreature}}
	if err := e.validate(); err != nil {
		t.Errorf("valid RemoveCounters should validate, got %v", err)
	}
	if got := e.Text(); got != "remove each growth counter from "+SelfName {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	tree := g.AddArtifact(NewCard("Vineapple Tree", Untamed, Artifact, Rare), 0)
	g.PlaceCounter(tree, CounterGrowth, 2)
	g.PlaceCounter(tree, CounterDoom, 1)
	e.Resolve(&EffectContext{Resolver: g, Source: tree, Controller: 0})
	if got := g.CountersOn(tree, CounterGrowth); got != 0 {
		t.Errorf("growth counters after remove = %d, want 0", got)
	}
	if got := g.CountersOn(tree, CounterDoom); got != 1 {
		t.Errorf("doom counters after remove = %d, want 1 (untouched)", got)
	}
	// Removing a kind the card no longer carries is a no-op.
	e.Resolve(&EffectContext{Resolver: g, Source: tree, Controller: 0})
	if got := g.CountersOn(tree, CounterGrowth); got != 0 {
		t.Errorf("growth counters after second remove = %d, want 0", got)
	}
}

// TestRemoveCountersAmount covers The Colosseum's partial removal: RemoveCounters
// with an Amount takes just that many, dropping the entry once it empties, and
// renders a counted noun. A negative Amount is rejected.
func TestRemoveCountersAmount(t *testing.T) {
	if err := (RemoveCounters{
		Kind:   CounterGlory,
		Target: Target{Kind: TargetThisCreature},
		Amount: -1,
	}).validate(); err == nil {
		t.Error("RemoveCounters with a negative Amount should not validate")
	}
	e := RemoveCounters{Kind: CounterGlory, Target: Target{Kind: TargetThisCreature}, Amount: 6}
	if got := e.Text(); got != "remove 6 glory counters from "+SelfName {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	arena := g.AddArtifact(NewCard("The Colosseum", Saurian, Artifact, Rare), 0)
	g.PlaceCounter(arena, CounterGlory, 8)
	e.Resolve(&EffectContext{Resolver: g, Source: arena, Controller: 0})
	if got := g.CountersOn(arena, CounterGlory); got != 2 {
		t.Errorf("glory counters after removing 6 of 8 = %d, want 2", got)
	}

	// Removing more than remain empties the entry; a non-positive n is a no-op.
	g.RemoveCountersN(arena, CounterGlory, 0)
	if got := g.CountersOn(arena, CounterGlory); got != 2 {
		t.Errorf("removing 0 changed the count to %d, want 2", got)
	}
	g.RemoveCountersN(arena, CounterGlory, 5)
	if got := g.CountersOn(arena, CounterGlory); got != 0 {
		t.Errorf("removing 5 of 2 = %d, want 0", got)
	}
	// Removing from a card that carries none is a no-op.
	g.RemoveCountersN(arena, CounterGlory, 1)
	if got := g.CountersOn(arena, CounterGlory); got != 0 {
		t.Errorf("removing from an empty entry = %d, want 0", got)
	}
}

// TestVineappleTreeGrowthCycle covers the whole card in the engine: the Action
// places a growth counter, each counter raises every player's key cost, and
// forging a key sheds them all.
func TestVineappleTreeGrowthCycle(t *testing.T) {
	tree := NewCard("Vineapple Tree", Untamed, Artifact, Rare,
		WithKeyCost(NewKeyCostChange(EachPlayer, 1).Per(CountersOnThis{Kind: CounterGrowth})),
		WithAbility(TriggerAfterPlayerForgesKey,
			RemoveCounters{Kind: CounterGrowth, Target: Target{Kind: TargetThisCreature}}),
		WithAbility(TriggerAction,
			PlaceCounter{Amount: 1, Kind: CounterGrowth, Target: Target{Kind: TargetThisCreature}}))

	g := NewGame("A", "B", 1)
	id := g.AddArtifact(tree, 0)

	if got := g.CurrentKeyCost(0); got != KeyCost {
		t.Errorf("key cost with no counters = %d, want %d", got, KeyCost)
	}

	// The Action places a growth counter on the tree.
	g.useActionOf(0, id)
	if got := g.CountersOn(id, CounterGrowth); got != 1 {
		t.Fatalf("growth counters after action = %d, want 1", got)
	}
	// One growth counter raises each player's key cost by one Æmber.
	if got := g.CurrentKeyCost(0); got != KeyCost+1 {
		t.Errorf("owner key cost with one counter = %d, want %d", got, KeyCost+1)
	}
	if got := g.CurrentKeyCost(1); got != KeyCost+1 {
		t.Errorf("opponent key cost with one counter = %d, want %d", got, KeyCost+1)
	}

	g.PlaceCounter(id, CounterGrowth, 2) // three total
	if got := g.CurrentKeyCost(0); got != KeyCost+3 {
		t.Errorf("key cost with three counters = %d, want %d", got, KeyCost+3)
	}

	// Forging a key sheds every growth counter, dropping the key cost back.
	g.forgeKeyFree(0)
	if got := g.CountersOn(id, CounterGrowth); got != 0 {
		t.Errorf("growth counters after forge = %d, want 0", got)
	}
	if got := g.CurrentKeyCost(0); got != KeyCost {
		t.Errorf("key cost after forge = %d, want %d", got, KeyCost)
	}
}

// TestPlaceCounterPer covers Book of Malefaction's reaction: a Per count scales
// the counters placed "for each Æmber stolen" and leads the rendered sentence.
func TestPlaceCounterPer(t *testing.T) {
	e := PlaceCounter{Amount: 1,
		Kind:   CounterWarrant,
		Target: Target{Kind: TargetThisCreature},
		Per:    AemberStolenThisEvent{},
	}
	if got := e.Text(); got != "for each Æmber stolen, put a warrant counter on {self}" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("book", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	ctx.Produced.AemberStolen = 3
	e.Resolve(ctx)
	if got := g.CountersOn(src, CounterWarrant); got != 3 {
		t.Errorf("warrant counters = %d, want 3", got)
	}
}

// TestRemoveCountersGate covers Book of Malefaction's Omni: RemoveCounters reports
// whether it took a counter off, so it can gate a Then, and renders singular text
// for Amount 1.
func TestRemoveCountersGate(t *testing.T) {
	e := RemoveCounters{
		Kind:   CounterWarrant,
		Target: Target{Kind: TargetThisCreature},
		Amount: 1,
	}
	if got := e.Text(); got != "remove a warrant counter from {self}" {
		t.Errorf("singular text = %q", got)
	}

	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("book", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	// With no counter the gate does nothing and reports false.
	if e.resolveGate(ctx) {
		t.Error("removing from a card with no counter should report false")
	}

	// With two counters it takes one off and reports true.
	g.PlaceCounter(src, CounterWarrant, 2)
	if !e.resolveGate(ctx) {
		t.Error("removing an existing counter should report true")
	}
	if got := g.CountersOn(src, CounterWarrant); got != 1 {
		t.Errorf("counters after removing one = %d, want 1", got)
	}
}
