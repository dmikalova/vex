package engine

import "testing"

// TestAemberPresenceIsTotal pins what each Æmber presence admits and the clause
// it prints, including the zero value, which asks nothing about Æmber and
// appends no clause.
func TestAemberPresenceIsTotal(t *testing.T) {
	want := map[AemberPresence]struct {
		bare   bool
		loaded bool
		clause string
	}{
		aemberAny:  {true, true, "each creature"},
		AemberSome: {false, true, "each creature with Æmber on it"},
		AemberNone: {true, false, "each creature with no Æmber on it"},
	}
	for _, presence := range append(allAemberPresences(), aemberAny) {
		expected, ok := want[presence]
		if !ok {
			t.Fatalf("Æmber presence %d has no pinned row; add one", presence)
		}
		if got := presence.admits(0); got != expected.bare {
			t.Errorf("presence %d admits(0) = %v, want %v", presence, got, expected.bare)
		}
		if got := presence.admits(1); got != expected.loaded {
			t.Errorf("presence %d admits(1) = %v, want %v", presence, got, expected.loaded)
		}
		if got := presence.clause("each creature"); got != expected.clause {
			t.Errorf("presence %d clause = %q, want %q", presence, got, expected.clause)
		}
		if got := presence.filters(); got != (presence != aemberAny) {
			t.Errorf("presence %d filters = %v", presence, got)
		}
	}
}
