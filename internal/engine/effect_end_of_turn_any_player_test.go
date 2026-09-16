package engine

import "testing"

// TestAfterAnyPlayerEndOfTurnResolvesAsActivePlayer covers the whole-board
// end-of-turn trigger (Pincerator): it fires at the end of every player's turn —
// its owner's and the opponent's — and resolves as the player whose turn is ending,
// so a GainAember{Controller} pays that active player rather than the artifact's
// controller.
func TestAfterAnyPlayerEndOfTurnResolvesAsActivePlayer(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// The artifact is player 0's, but its ability pays whoever's turn is ending.
	g.AddArtifact(NewCard("pincer", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterAnyPlayerEndOfTurn,
			GainAember{Player: Controller, Amount: 1})), 0)

	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.State.Aember[0] != 1 {
		t.Fatalf("after player 0's turn end: aember[0] = %d, want 1", g.State.Aember[0])
	}
	if g.State.Aember[1] != 0 {
		t.Fatalf("after player 0's turn end: aember[1] = %d, want 0", g.State.Aember[1])
	}

	g.StartTurn(1)
	g.EndPlayPhase(1)
	if g.State.Aember[1] != 1 {
		t.Fatalf(
			"after player 1's turn end: aember[1] = %d, want 1 (resolves as the active player)",
			g.State.Aember[1],
		)
	}
	if g.State.Aember[0] != 1 {
		t.Fatalf(
			"after player 1's turn end: aember[0] = %d, want 1 (unchanged)",
			g.State.Aember[0],
		)
	}
}

// TestAfterAnyPlayerEndOfTurnFiresForBothPlayersCards confirms the window scans
// both battlelines, not just the active player's, so an opponent-owned end-of-turn
// artifact still fires on the active player's turn.
func TestAfterAnyPlayerEndOfTurnFiresForBothPlayersCards(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// Player 1 owns the artifact; it still fires at the end of player 0's turn and
	// pays the active player (player 0).
	g.AddArtifact(NewCard("pincer", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterAnyPlayerEndOfTurn,
			GainAember{Player: Controller, Amount: 1})), 1)

	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.State.Aember[0] != 1 {
		t.Fatalf(
			"aember[0] = %d, want 1 (an opponent-owned artifact fires and pays the active player)",
			g.State.Aember[0],
		)
	}
	if g.State.Aember[1] != 0 {
		t.Fatalf("aember[1] = %d, want 0", g.State.Aember[1])
	}
}

// TestAfterAnyPlayerEndOfTurnPrefix covers the printed prefix, mirroring the
// start-of-turn whole-board trigger.
func TestAfterAnyPlayerEndOfTurnPrefix(t *testing.T) {
	got, _ := TriggerAfterAnyPlayerEndOfTurn.prefix()
	if want := "At the end of each player's turn, "; got != want {
		t.Errorf("prefix = %q, want %q", got, want)
	}
}
