package engine

// This file holds the projection seam (ADR 0039): the single point where the
// authoritative GameState is turned into what a particular viewer is allowed to
// see. Today the projection is the identity — every viewer sees the whole state —
// but routing all client rendering through it means a future networked server can
// redact an opponent's hidden zones (hand, deck order, facedown cards) per viewer
// without changing any caller. Clients render from the View, never from raw state.

// View is what one viewer sees of a match: the projected state plus the viewer it
// was projected for. Today State is the whole GameState unredacted; the wrapper
// exists so redaction can hollow out hidden zones later without changing the type
// clients consume.
type View struct {
	// Viewer is the player index the projection was made for.
	Viewer int
	// State is the projected game state this viewer may see.
	State GameState
}

// Project returns the View that viewer is allowed to see of state. It is the
// identity projection today (the whole state, unredacted); the seam is here so
// server-side redaction can be added later at this one point.
func Project(state GameState, viewer int) View {
	return View{
		Viewer: viewer,
		State:  state,
	}
}
