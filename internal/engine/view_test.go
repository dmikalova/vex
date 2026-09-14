package engine

import "testing"

// The projection is the identity today: a viewer sees the whole state unchanged,
// and the View records which viewer it was made for. The seam exists so redaction
// can be added here later without changing callers.
func TestProjectIsIdentity(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 3
	g.State.Aember[1] = 5

	v := Project(g.State, 1)
	if v.Viewer != 1 {
		t.Errorf("Viewer = %d, want 1", v.Viewer)
	}
	if v.State != g.State {
		t.Error("the identity projection changed the state")
	}
}
