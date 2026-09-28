package engine

import "testing"

// TestPositionTextIsTotal pins how each battleline place prints: the flank
// renders as an adjective on the noun, the rest as a trailing clause, and the
// zero value prints neither. A place renders one way or the other, never both,
// so a newly added value cannot silently print nothing.
func TestPositionTextIsTotal(t *testing.T) {
	want := map[Position]struct {
		adjective string
		clause    string
	}{
		positionAny:           {"", "each creature"},
		PositionOnFlank:       {"flank", "each creature"},
		PositionNotOnFlank:    {"", "each creature that is not on a flank"},
		PositionCenter:        {"", "each creature in the center of its controller's battleline"},
		PositionRightOfSource: {"", "each creature to the right of " + SelfName},
		PositionLeftOfSource:  {"", "each creature to the left of " + SelfName},
	}
	for _, position := range append(allPositions(), positionAny) {
		expected, ok := want[position]
		if !ok {
			t.Fatalf("position %d has no pinned row; add one", position)
		}
		if got := position.adjective(); got != expected.adjective {
			t.Errorf("position %d adjective = %q, want %q", position, got, expected.adjective)
		}
		if got := position.clause("each creature"); got != expected.clause {
			t.Errorf("position %d clause = %q, want %q", position, got, expected.clause)
		}
		if position.filters() && expected.adjective == "" && expected.clause == "each creature" {
			t.Errorf("position %d narrows the set but prints nothing", position)
		}
		if got := position.filters(); got != (position != positionAny) {
			t.Errorf("position %d filters = %v", position, got)
		}
	}
}

// TestZeroPositionAdmitsEveryCard checks that a target asking nothing about
// position admits a card wherever it stands — the branch a filtered target with
// no position of its own takes.
func TestZeroPositionAdmitsEveryCard(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 3), 0)
	middle := g.AddToBattleline(testCreature("middle", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	for _, id := range []LocalID{left, middle} {
		if !positionAny.admits(ctx, id) {
			t.Errorf("the zero Position rejected %d", id)
		}
	}
}
