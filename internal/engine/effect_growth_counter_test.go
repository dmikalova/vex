package engine

import "testing"

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
			PlaceCounter{Kind: CounterGrowth, Target: Target{Kind: TargetThisCreature}}))

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
