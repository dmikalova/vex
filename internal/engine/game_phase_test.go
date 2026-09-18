package engine

import "testing"

func TestPhaseNames(t *testing.T) {
	for _, tc := range []struct {
		phase Phase
		want  string
		step  string
	}{
		{PhaseStartOfTurn, "start of turn", "1. Start of turn"},
		{PhaseForge, "forge a key", "2. Forge a key"},
		{PhaseChooseHouse, "choose a house", "3. Choose a house"},
		{PhaseArchives, "archives", "4. Archives"},
		{PhasePlay, "main", "5. Main phase"},
		{PhaseReady, "ready", "6. Ready"},
		{PhaseDraw, "draw", "7. Draw"},
		{PhaseEndOfTurn, "end of turn", "8. End of turn"},
		{phaseUnset, "no phase", ""},
	} {
		if got := tc.phase.String(); got != tc.want {
			t.Errorf("Phase(%d).String() = %q, want %q", tc.phase, got, tc.want)
		}
		if got := tc.phase.rulebookStep(); got != tc.step {
			t.Errorf("Phase(%d).rulebookStep() = %q, want %q", tc.phase, got, tc.step)
		}
	}
	if phaseUnset.valid() {
		t.Error("the zero phase must be invalid")
	}
	if !PhasePlay.valid() {
		t.Error("PhasePlay must be valid")
	}
}

func TestTurnWalksThroughItsPhases(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	if g.State.Phase.valid() {
		t.Errorf("a game that has not begun a turn is in phase %v, want none", g.State.Phase)
	}

	g.StartTurn(0)
	if g.State.Phase != PhaseChooseHouse {
		t.Errorf("after StartTurn the phase is %v, want choose a house", g.State.Phase)
	}
	// A frontend reads the phase through the accessor, not the state field.
	if g.Phase() != PhaseChooseHouse {
		t.Errorf("Phase() = %v, want choose a house", g.Phase())
	}

	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if g.State.Phase != PhasePlay {
		t.Errorf("after ChooseHouse the phase is %v, want play", g.State.Phase)
	}

	g.EndPlayPhase(0)
	if g.State.Phase != PhaseEndOfTurn {
		t.Errorf("after EndPlayPhase the phase is %v, want end of turn", g.State.Phase)
	}
}

func TestEndPhaseSkipsAnOpenPhase(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.StartTurn(0)
	// The choose-a-house phase normally blocks for the player. Ending it early
	// lets the loop walk on to the next phase that blocks.
	g.EndPhase()
	g.runPhases()
	if g.State.Phase != PhasePlay {
		t.Errorf("phase = %v, want play (the loop should skip the ended phase)", g.State.Phase)
	}
	if g.State.PhaseEnded {
		t.Error("entering a phase must clear the early-end flag")
	}
}

func TestConcedeHandsTheGameToTheOpponent(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)

	g.Concede(0)

	if g.Winner() != 1 {
		t.Fatalf("winner = %d, want 1 (opponent of the conceding player)", g.Winner())
	}
	last := g.Log[len(g.Log)-1].Entry
	if _, ok := last.(PlayerConceded); !ok {
		t.Errorf("last log entry = %#v, want PlayerConceded", last)
	}

	// A second concede is a no-op: the game is already decided.
	before := len(g.Log)
	g.Concede(1)
	if g.Winner() != 1 || len(g.Log) != before {
		t.Errorf("conceding a decided game changed it: winner %d, log grew by %d",
			g.Winner(), len(g.Log)-before)
	}
}

func TestNoPhaseLogAfterGameWon(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.State.ForgeCanonicalKeys(0, KeysToWin-1)
	g.State.Aember[0] = KeyCost

	// Forging the third key wins the game mid-forge-phase, before the loop would
	// otherwise walk on to the archives phase.
	g.StartTurn(0)

	if g.Winner() != 0 {
		t.Fatalf("winner = %d, want 0", g.Winner())
	}
	last := g.Log[len(g.Log)-1].Entry
	if _, ok := last.(GameWon); !ok {
		t.Errorf("last log entry = %#v, want GameWon (no phase entered after the win)", last)
	}
}

func TestStartOfTurnAbilitiesResolveBeforeForging(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.AddArtifact(NewCard("dawn", Brobnar, Artifact, Rare,
		WithAbility(TriggerStartOfTurn, GainAember{Player: Controller, Amount: 6})), 0)

	g.StartTurn(0)

	// The Æmber arrived in time to pay for the turn's forge.
	if g.Keys(0) != 1 {
		t.Errorf(
			"keys = %d, want 1 (start-of-turn Æmber should pay for the forge)",
			g.Keys(0),
		)
	}
}

func TestEndOfTurnAbilitiesResolveAfterReadyAndDraw(t *testing.T) {
	g := started(t)
	var handAtTrigger, exhaustedAtTrigger int
	watcher := gameEffect{fn: func() {
		handAtTrigger = len(g.Hand(0))
		for _, id := range g.allInPlay(0) {
			if g.State.Cards[id].Exhausted {
				exhaustedAtTrigger++
			}
		}
	}}
	c := g.AddToBattleline(NewCard("watcher", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerEndOfTurn, watcher)), 0)
	g.State.Cards[c].Exhausted = true
	for i := 0; i < HandSize; i++ {
		g.AddToDeck(testCreature("stock", 1), 0)
	}

	g.EndPlayPhase(0)

	if exhaustedAtTrigger != 0 {
		t.Errorf(
			"%d cards were still exhausted when the end-of-turn ability resolved; ready runs first",
			exhaustedAtTrigger,
		)
	}
	if handAtTrigger != HandSize {
		t.Errorf("hand was %d when the end-of-turn ability resolved, want %d (draw runs first)",
			handAtTrigger, HandSize)
	}
}

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

func TestEndOfTurnTriggerFires(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.StartTurn(0)
	g.State.Aember[1] = 3
	g.AddToBattleline(NewCard("Shaffles", Dis, Creature, Common, WithPower(2),
		WithAbility(TriggerEndOfTurn, LoseAember{Player: Opponent, Amount: 1})), 0)

	g.EndPlayPhase(0)

	if g.State.Aember[1] != 2 {
		t.Errorf("opponent Æmber = %d, want 2 after end-of-turn drain", g.State.Aember[1])
	}
}
