package engine

import "testing"

// This test covers the Omega keyword: a card with Omega ends the current step of
// the turn the moment it resolves.

func TestOmegaEndsStep(t *testing.T) {
	g := started(t)
	omega := g.AddToHand(
		NewCard("Gateway", Brobnar, Tactic, Common, WithKeywords(Omega)),
		0,
	)
	// A second card left in hand would be playable if the step had not ended.
	g.AddToHand(testCreature("leftover", 3), 0)
	if err := g.PlayAction(0, handIdxByID(g, 0, omega)); err != nil {
		t.Fatalf("play Omega action: %v", err)
	}
	if g.State.Phase != PhaseEndOfTurn {
		t.Errorf("phase after Omega = %v, want PhaseEndOfTurn", g.State.Phase)
	}
}
