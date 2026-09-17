package engine

import "testing"

// TestItAttachedToThisOrNeighborCondText covers the condition Commander Dhrxgar
// reads, naming the source card via the self placeholder.
func TestItAttachedToThisOrNeighborCondText(t *testing.T) {
	want := "if it is attached to " + SelfName + " or one of its neighbors"
	if got := (ItAttachedToThisOrNeighbor{}).CondText(); got != want {
		t.Errorf("CondText = %q, want %q", got, want)
	}
}
