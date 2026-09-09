package engine

import "testing"

// TestAfterAnyPlayerStartOfTurnResolvesAsActivePlayer covers the whole-board
// start-of-turn trigger (Gambling Den, General Order 24): it fires at the start of
// every player's turn — its owner's and the opponent's — and resolves as the player
// whose turn is starting, so a GainAember{Controller} pays that active player rather
// than the artifact's controller.
func TestAfterAnyPlayerStartOfTurnResolvesAsActivePlayer(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// The artifact is player 0's, but its ability pays whoever's turn is starting.
	g.AddArtifact(NewCard("den", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterAnyPlayerStartOfTurn,
			GainAember{Player: Controller, Amount: 1})), 0)

	g.StartTurn(0)
	if g.State.Aember[0] != 1 {
		t.Fatalf("after player 0's turn start: aember[0] = %d, want 1", g.State.Aember[0])
	}
	if g.State.Aember[1] != 0 {
		t.Fatalf("after player 0's turn start: aember[1] = %d, want 0", g.State.Aember[1])
	}

	g.StartTurn(1)
	if g.State.Aember[1] != 1 {
		t.Fatalf(
			"after player 1's turn start: aember[1] = %d, want 1 (resolves as the active player)",
			g.State.Aember[1],
		)
	}
	if g.State.Aember[0] != 1 {
		t.Fatalf(
			"after player 1's turn start: aember[0] = %d, want 1 (unchanged)",
			g.State.Aember[0],
		)
	}
}

// TestAfterAnyPlayerStartOfTurnFiresForBothPlayersCards confirms the window scans
// both battlelines' artifacts, not just the active player's, so an opponent-owned
// start-of-turn artifact still fires on the active player's turn.
func TestAfterAnyPlayerStartOfTurnFiresForBothPlayersCards(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// Player 1 owns the artifact; it still fires at the start of player 0's turn and
	// pays the active player (player 0).
	g.AddArtifact(NewCard("den", Brobnar, Artifact, Rare,
		WithAbility(TriggerAfterAnyPlayerStartOfTurn,
			GainAember{Player: Controller, Amount: 1})), 1)

	g.StartTurn(0)
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
