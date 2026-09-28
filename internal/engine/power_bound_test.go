package engine

import "testing"

// TestPowerBoundAdmitsIsTotal pins what each bound admits, including the zero
// value, which bounds nothing and so admits every power. A power filter is one
// field rather than a bound beside a parity flag, so the zero value is the only
// way to say "no power question" and every other kind must answer one.
func TestPowerBoundAdmitsIsTotal(t *testing.T) {
	powers := []int{0, 1, 2, 3, 4}
	want := map[PowerBoundKind][]bool{
		boundUnset:   {true, true, true, true, true},
		BoundAtMost:  {true, true, true, true, false},
		BoundAtLeast: {false, false, false, true, true},
		BoundExactly: {false, false, false, true, false},
		BoundOdd:     {false, true, false, true, false},
		BoundEven:    {true, false, true, false, true},
		// The source creature registered below has power 3, so a lower power passes.
		BoundLessThanSource: {true, true, true, false, false},
	}
	g := NewGame("A", "B", 1)
	source := g.AddToBattleline(testCreature("source", 3), 0)
	ctx := newTestContext(g, source)
	for kind, admits := range want {
		bound := PowerBound{Kind: kind, Amount: 3}
		for i, power := range powers {
			if got := bound.admits(ctx, power); got != admits[i] {
				t.Errorf("PowerBound{%d, 3}.admits(%d) = %v, want %v",
					kind, power, got, admits[i])
			}
		}
	}
	for _, kind := range allPowerBoundKinds() {
		if _, ok := want[kind]; !ok {
			t.Errorf("bound kind %d has no pinned admits row; add one", kind)
		}
		if !(PowerBound{Kind: kind}).filters() {
			t.Errorf("bound kind %d does not report that it filters", kind)
		}
	}
	if (PowerBound{}).filters() {
		t.Error("the zero PowerBound reports that it filters")
	}
}

// TestPowerBoundClauseIsTotal pins the printed clause of every bound kind, so a
// newly added kind cannot fall through clause's default and print nothing where
// a card says something.
func TestPowerBoundClauseIsTotal(t *testing.T) {
	want := map[PowerBoundKind]string{
		boundUnset:   "each creature",
		BoundAtMost:  "each creature with power 3 or lower",
		BoundAtLeast: "each creature with power 3 or higher",
		BoundExactly: "each creature with power 3",
		BoundOdd:     "each creature with odd power",
		BoundEven:    "each creature with even power",
		// A source-relative bound names the card rather than a number.
		BoundLessThanSource: "each creature with lower power than " + SelfName,
	}
	for _, kind := range append(allPowerBoundKinds(), boundUnset) {
		expected, ok := want[kind]
		if !ok {
			t.Fatalf("bound kind %d has no pinned clause; add one", kind)
		}
		got := PowerBound{Kind: kind, Amount: 3}.clause("each creature")
		if got != expected {
			t.Errorf("bound kind %d clause = %q, want %q", kind, got, expected)
		}
	}
}

// newTestContext is the minimal resolution context the bound tests read: a game,
// its controller, and the source card whose power a source-relative bound
// measures against.
func newTestContext(g *Game, source LocalID) *EffectContext {
	return &EffectContext{Resolver: g, Controller: 0, Source: source}
}
