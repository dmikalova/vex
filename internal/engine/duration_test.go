package engine

import "testing"

// TestDurationsAreRealAndNamed checks that Durations lists only real durations
// (never the unset zero or the count sentinel) and that each names itself with a
// distinct non-empty label.
func TestDurationsAreRealAndNamed(t *testing.T) {
	seen := map[string]bool{}
	for _, d := range Durations() {
		if !d.valid() {
			t.Errorf("Durations included the unset value %d", d)
		}
		if d >= durationCount {
			t.Errorf("Durations included the count sentinel %d", d)
		}
		name := d.String()
		if name == "" {
			t.Errorf("duration %d has an empty String", d)
		}
		if seen[name] {
			t.Errorf("duration name %q is not distinct", name)
		}
		seen[name] = true
	}
	if len(Durations()) != int(durationCount)-1 {
		t.Errorf("Durations len = %d, want %d", len(Durations()), int(durationCount)-1)
	}
}

// TestDurationStringDefaultsEmpty checks that a non-real duration renders empty,
// covering the default branch of String.
func TestDurationStringDefaultsEmpty(t *testing.T) {
	if got := durationUnset.String(); got != "" {
		t.Errorf("durationUnset.String() = %q, want empty", got)
	}
	if got := durationCount.String(); got != "" {
		t.Errorf("durationCount.String() = %q, want empty", got)
	}
}

// TestDurationClauseIsTotal pins the standard turn-window clause of every real
// duration, so the single-sourced phrasing cannot drift and a newly added duration
// cannot silently fall through durationClause's default to the wrong phrase — a new
// value not listed here fails the lookup. Forever carries no standing clause and
// renders empty. The source token is threaded through so UntilThisLeavesPlay names
// the card whose leaving ends the effect.
func TestDurationClauseIsTotal(t *testing.T) {
	want := map[Duration]string{
		RemainderOfPlayerTurn: "for the remainder of the turn",
		OpponentNextTurn:      "during your opponent's next turn",
		StartOfPlayerNextTurn: "until the start of your next turn",
		EndOfPlayerNextTurn:   "until the end of your next turn",
		UntilThisLeavesPlay:   "until {self} leaves play",
		Forever:               "",
	}
	for _, d := range Durations() {
		expected, ok := want[d]
		if !ok {
			t.Fatalf("durationClause: no expected clause pinned for %v; add its "+
				"clause to durationClause and pin it here", d)
		}
		if got := durationClause(d, "{self}"); got != expected {
			t.Errorf("durationClause(%v) = %q, want %q", d, got, expected)
		}
	}
	if got := durationClause(durationUnset, "{self}"); got != "" {
		t.Errorf("durationClause(durationUnset) = %q, want empty", got)
	}
}

// TestWindowClauseIsTotal pins the subject-framed window clause of every real
// duration, threading a sample possessive so the next-turn family renders relative
// to the effect's own subject ("their next turn") rather than the controller's
// absolute frame. Like durationClause, it must be total — a newly added duration
// not pinned here fails the lookup, so it cannot fall through windowClause's default
// to the wrong phrase. Forever renders empty; UntilThisLeavesPlay names the card
// whose leaving ends the effect.
func TestWindowClauseIsTotal(t *testing.T) {
	want := map[Duration]string{
		RemainderOfPlayerTurn: "for the remainder of the turn",
		OpponentNextTurn:      "during their next turn",
		StartOfPlayerNextTurn: "until the start of their next turn",
		EndOfPlayerNextTurn:   "until the end of their next turn",
		UntilThisLeavesPlay:   "until {self} leaves play",
		Forever:               "",
	}
	for _, d := range Durations() {
		expected, ok := want[d]
		if !ok {
			t.Fatalf("windowClause: no expected clause pinned for %v; add its "+
				"clause to windowClause and pin it here", d)
		}
		if got := windowClause(d, "their", "{self}"); got != expected {
			t.Errorf("windowClause(%v) = %q, want %q", d, got, expected)
		}
	}
	if got := windowClause(durationUnset, "their", "{self}"); got != "" {
		t.Errorf("windowClause(durationUnset) = %q, want empty", got)
	}
}
