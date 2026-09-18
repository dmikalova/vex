package engine

import "testing"

// TestPromptKindsAreTheCatalogMinusUnset checks the exported catalog lists every
// real kind and never the invalid zero.
func TestPromptKindsAreTheCatalogMinusUnset(t *testing.T) {
	kinds := PromptKinds()
	if len(kinds) != 7 {
		t.Fatalf("PromptKinds returned %d kinds, want 7", len(kinds))
	}
	for _, k := range kinds {
		if k == promptKindUnset {
			t.Errorf("PromptKinds included the invalid zero value")
		}
	}
}

// TestPromptKindString names every kind, including the invalid zero and an
// out-of-range value.
func TestPromptKindString(t *testing.T) {
	cases := []struct {
		kind PromptKind
		want string
	}{
		{PromptCreature, "creature"},
		{PromptCardOrDecline, "card-or-decline"},
		{PromptOption, "option"},
		{PromptPosition, "position"},
		{PromptReaction, "reaction"},
		{PromptOrder, "order"},
		{PromptBadge, "badge"},
		{promptKindUnset, "unset"},
		{PromptKind(99), "invalid"},
	}
	for _, tc := range cases {
		if got := tc.kind.String(); got != tc.want {
			t.Errorf("PromptKind(%d).String() = %q, want %q", tc.kind, got, tc.want)
		}
	}
}
