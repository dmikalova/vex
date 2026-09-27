package engine

import (
	"reflect"
	"slices"
	"testing"
)

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
		WithAbility(TriggerStartOfTurn, GainAember{
			Player: Controller,
			Amount: 6,
		})), 0)

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
		for _, id := range g.creaturesAndArtifacts(0) {
			if g.State.Cards[id].Exhausted {
				exhaustedAtTrigger++
			}
		}
	}}
	c := g.AddToBattleline(NewCard("watcher", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerEndOfTurn, watcher)), 0)
	g.State.Cards[c].Exhausted = true
	for range HandSize {
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

// TestEachPlayerEndOfTurnResolvesAsActivePlayer covers the whole-board
// end-of-turn trigger (Pincerator): it fires at the end of every player's turn —
// its owner's and the opponent's — and resolves as the player whose turn is ending,
// so a GainAember{Controller} pays that active player rather than the artifact's
// controller.
func TestEachPlayerEndOfTurnResolvesAsActivePlayer(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// The artifact is player 0's, but its ability pays whoever's turn is ending.
	g.AddArtifact(NewCard("pincer", Brobnar, Artifact, Rare,
		WithEachPlayerAbility(TriggerEndOfTurn,
			GainAember{
				Player: Controller,
				Amount: 1,
			})), 0)

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

// TestEachPlayerEndOfTurnFiresForBothPlayersCards confirms the window scans
// both battlelines, not just the active player's, so an opponent-owned end-of-turn
// artifact still fires on the active player's turn.
func TestEachPlayerEndOfTurnFiresForBothPlayersCards(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// Player 1 owns the artifact; it still fires at the end of player 0's turn and
	// pays the active player (player 0).
	g.AddArtifact(NewCard("pincer", Brobnar, Artifact, Rare,
		WithEachPlayerAbility(TriggerEndOfTurn,
			GainAember{
				Player: Controller,
				Amount: 1,
			})), 1)

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

// TestEachPlayerEndOfTurnPrefix covers the printed prefix, mirroring the
// start-of-turn whole-board trigger.
func TestEachPlayerEndOfTurnPrefix(t *testing.T) {
	got, _ := abilityPrefix(Ability{
		Trigger:    TriggerEndOfTurn,
		EachPlayer: true,
	})
	if want := "At the end of each player's turn, "; got != want {
		t.Errorf("prefix = %q, want %q", got, want)
	}
}

// TestEachPlayerStartOfTurnResolvesAsActivePlayer covers the whole-board
// start-of-turn trigger (Gambling Den, General Order 24): it fires at the start of
// every player's turn — its owner's and the opponent's — and resolves as the player
// whose turn is starting, so a GainAember{Controller} pays that active player rather
// than the artifact's controller.
func TestEachPlayerStartOfTurnResolvesAsActivePlayer(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// The artifact is player 0's, but its ability pays whoever's turn is starting.
	g.AddArtifact(NewCard("den", Brobnar, Artifact, Rare,
		WithEachPlayerAbility(TriggerStartOfTurn,
			GainAember{
				Player: Controller,
				Amount: 1,
			})), 0)

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

// TestEachPlayerStartOfTurnFiresForBothPlayersCards confirms the window scans
// both battlelines' artifacts, not just the active player's, so an opponent-owned
// start-of-turn artifact still fires on the active player's turn.
func TestEachPlayerStartOfTurnFiresForBothPlayersCards(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// Player 1 owns the artifact; it still fires at the start of player 0's turn and
	// pays the active player (player 0).
	g.AddArtifact(NewCard("den", Brobnar, Artifact, Rare,
		WithEachPlayerAbility(TriggerStartOfTurn,
			GainAember{
				Player: Controller,
				Amount: 1,
			})), 1)

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

// TestEachPlayerStartOfTurnPrefix covers the printed prefix, mirroring the
// end-of-turn whole-board trigger.
func TestEachPlayerStartOfTurnPrefix(t *testing.T) {
	got, _ := abilityPrefix(Ability{
		Trigger:    TriggerStartOfTurn,
		EachPlayer: true,
	})
	if want := "At the start of each player's turn, "; got != want {
		t.Errorf("prefix = %q, want %q", got, want)
	}
}

// TestOwnScopeEndOfTurnDoesNotFireOnOpponentTurn pins the skip in the whole-board
// scan: an ordinary (non-EachPlayer) end-of-turn ability on the opponent's board
// watches only its own controller's turn, so it stays silent when the active
// player ends their turn.
func TestOwnScopeEndOfTurnDoesNotFireOnOpponentTurn(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.State.Aember[1] = 3
	// Player 1's own-scope end-of-turn drain must not fire at the end of player 0's
	// turn — it is not EachPlayer-scoped, so the whole-board scan skips it.
	g.AddArtifact(NewCard("drain", Brobnar, Artifact, Rare,
		WithAbility(TriggerEndOfTurn, LoseAember{
			Player: Opponent,
			Amount: 1,
		})), 1)

	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.State.Aember[0] != 0 {
		t.Fatalf(
			"aember[0] = %d, want 0 (opponent's own-scope drain must not fire)",
			g.State.Aember[0],
		)
	}
}

func TestEndOfTurnTriggerFires(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.StartTurn(0)
	g.State.Aember[1] = 3
	g.AddToBattleline(NewCard("Shaffles", Dis, Creature, Common, WithPower(2),
		WithAbility(TriggerEndOfTurn, LoseAember{
			Player: Opponent,
			Amount: 1,
		})), 0)

	g.EndPlayPhase(0)

	if g.State.Aember[1] != 2 {
		t.Errorf("opponent Æmber = %d, want 2 after end-of-turn drain", g.State.Aember[1])
	}
}

// TestAnimatorRevertsAfterEndOfTurnAbilities is the Animator plus Fangtooth
// Cavern ruling (docs/keyforge-master-rulebook.md): a card animated for the
// remainder of the turn is "a lasting effect that expires at the end of the turn,
// which is after Fangtooth Cavern's 'end of turn' effect resolves". So an
// end-of-turn ability that reads the board still sees the animated artifact as a
// creature, and the revert lands in the cleanup tail after that ability has
// resolved (ADR 0047).
func TestAnimatorRevertsAfterEndOfTurnAbilities(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	art := g.AddArtifact(testArtifact("animated"), 0)
	// Animator gives the artifact three +1 power counters before animating it, so
	// it is a 3-power creature and survives the settles along the way.
	g.AddPowerCounter(art, 3)
	var typeAtTrigger CardType
	var inBattlelineAtTrigger bool
	g.AddToBattleline(NewCard("cavern", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerEndOfTurn, gameEffect{fn: func() {
			typeAtTrigger = g.TypeOf(art)
			inBattlelineAtTrigger = slices.Contains(g.Battleline(0), art)
		}})), 0)
	g.StartTurn(0)
	TurnIntoCreature{
		Target:   Target{Kind: TargetThisCreature},
		Duration: RemainderOfPlayerTurn,
	}.Resolve(&EffectContext{
		Resolver:   g,
		Source:     art,
		Controller: 0,
	})

	g.EndPlayPhase(0)

	if typeAtTrigger != Creature {
		t.Errorf(
			"the end-of-turn ability saw the animated card as %v, want %v (the revert comes after)",
			typeAtTrigger, Creature,
		)
	}
	if !inBattlelineAtTrigger {
		t.Error("the end-of-turn ability should still see the animated card in the battleline")
	}
	if got := g.TypeOf(art); got != Artifact {
		t.Errorf("after the turn ended the card is %v, want %v", got, Artifact)
	}
}

// ---- witness tests for the end-of-turn cleanup tail ----
//
// The cleanup tail (expireTurnScoped, run after the end-of-turn abilities under
// ADR 0047) writes state directly rather than through the Resolver port, so it is
// outside the narration audit by design: refreshing armor, lifting the turn's
// bars, and rolling TurnHistory are not capabilities any card has. These tests
// cover it instead, one step each — what the step changed, and what it did or did
// not say about it.
//
// The silent steps are silent on purpose, and each test below says why. The rule
// they share: the turn boundary is not an event a card causes, and every
// "remainder of the turn" effect ends at it, so narrating each expiry would
// append a line per creature per turn for something the log already told the
// player was temporary. The expiries that ARE narrated are the ones that change
// what a card is (a reverting animation), not the ones that merely return it to
// its printed state.

// entriesAppended runs step and returns the type name of each log entry it
// appended, in order — the witness one phase step leaves in the log.
func entriesAppended(g *Game, step func()) []string {
	before := len(g.Log)
	step()
	var out []string
	for _, r := range g.Log[before:] {
		out = append(out, reflect.TypeOf(r.Entry).Name())
	}
	return out
}

// TestReadyStepNarratesTheCardsItReadied is the witness for readying: it names
// the cards that were actually turned upright, in one line, so a player can see
// what came back.
func TestReadyStepNarratesTheCardsItReadied(t *testing.T) {
	g := started(t)
	exhausted := g.AddToBattleline(testCreature("exhausted", 3), 0)
	g.AddToBattleline(testCreature("ready", 3), 0)
	g.State.Cards[exhausted].Exhausted = true

	got := entriesAppended(g, func() { g.readyPhase(0) })

	if !slices.Contains(got, "CardsReadied") {
		t.Fatalf("ready step appended %v, want a CardsReadied entry", got)
	}
	for _, r := range g.Log {
		if e, ok := r.Entry.(CardsReadied); ok {
			if !slices.Equal(e.Cards, []LocalID{exhausted}) {
				t.Errorf("CardsReadied names %v, want only the exhausted card %d",
					e.Cards, exhausted)
			}
		}
	}
	if g.State.Cards[exhausted].Exhausted {
		t.Error("the card is still exhausted after the ready step")
	}
}

// TestCleanupNarratesATemporaryAnimationReverting is the witness for the Animator
// revert. This expiry IS narrated: the card stops being a creature and leaves the
// battleline, which changes what the board is, so RevertedToArtifact tells the
// player why a creature they could see is gone.
func TestCleanupNarratesATemporaryAnimationReverting(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	art := g.AddArtifact(testArtifact("animated"), 0)
	g.AddPowerCounter(art, 3)
	g.StartTurn(0)
	TurnIntoCreature{
		Target:   Target{Kind: TargetThisCreature},
		Duration: RemainderOfPlayerTurn,
	}.Resolve(&EffectContext{
		Resolver:   g,
		Source:     art,
		Controller: 0,
	})

	got := entriesAppended(g, func() { g.expireTurnScoped(0) })

	if !slices.Contains(got, "RevertedToArtifact") {
		t.Fatalf("cleanup appended %v, want a RevertedToArtifact entry", got)
	}
	if typ := g.TypeOf(art); typ != Artifact {
		t.Errorf("the animated card is %v after cleanup, want %v", typ, Artifact)
	}
}

// TestCleanupLiftsTurnBarsSilently pins the decision that a bar expiring is not
// narrated. Arming one is not narrated either (ADR 0011): a bar narrates when it
// bites, as CardCannotBeUsed on the use it refuses. A line for the lift would be
// the only mention of a bar the log never announced.
func TestCleanupLiftsTurnBarsSilently(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(testCreature("fogbank", 3), 0)
	g.State.CannotFight[0] = Bar[bool]{
		Value:  true,
		Source: source,
	}

	got := entriesAppended(g, func() { g.expireTurnScoped(0) })

	if len(got) != 0 {
		t.Errorf("lifting the fight bar appended %v, want nothing", got)
	}
	if g.State.CannotFight[0].Value {
		t.Error("the fight bar is still in force after cleanup")
	}
}

// TestCleanupExpiresTemporaryBuffsSilently pins the decision that a "remainder of
// the turn" stat or keyword grant expiring is not narrated. The grant itself was
// narrated (CreatureGainedStats, CreatureGainedKeyword) and every such grant ends
// at the same moment, so a line per buffed creature per turn would restate the
// turn boundary rather than tell the player anything new.
func TestCleanupExpiresTemporaryBuffsSilently(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("buffed", 3), 0)
	g.GainStats(id, 2, 0)
	g.GrantKeyword(id, Skirmish)

	got := entriesAppended(g, func() { g.expireTurnScoped(0) })

	if len(got) != 0 {
		t.Errorf("expiring the turn's buffs appended %v, want nothing", got)
	}
	if bonus := g.State.Cards[id].TempPowerBonus; bonus != 0 {
		t.Errorf("temporary power bonus is %d after cleanup, want 0", bonus)
	}
	if g.State.Cards[id].GrantedKeywords != 0 {
		t.Error("granted keywords survived the cleanup")
	}
}

// TestCleanupRefreshesArmorSilently pins the decision that the armor refresh is
// not narrated. Armor spent (ArmorAbsorbed) and armor stripped (ArmorLost) are
// both narrated when they happen; returning a creature to its printed armor for
// the turn to come is the rule, not an outcome, and applies to every creature in
// play every turn.
func TestCleanupRefreshesArmorSilently(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("armored", 3, WithArmor(2)), 0)
	g.StripArmor(id)

	got := entriesAppended(g, func() { g.expireTurnScoped(0) })

	if len(got) != 0 {
		t.Errorf("refreshing armor appended %v, want nothing", got)
	}
	core := g.State.Cards[id]
	if core.ArmorRemaining != 2 || core.ArmorStripped != 0 {
		t.Errorf("armor after cleanup is %d remaining, %d stripped; want 2 and 0",
			core.ArmorRemaining, core.ArmorStripped)
	}
}

// TestCleanupRollsTurnHistorySilently pins the decision that the TurnHistory roll
// is not narrated: it is pure bookkeeping. The tallies it rolls record things the
// log already narrated as they happened (a key forged, a creature destroyed), and
// no player can see the tallies themselves — they exist only for the cards that
// ask "did your opponent forge a key on their previous turn?".
func TestCleanupRollsTurnHistorySilently(t *testing.T) {
	g := started(t)
	g.State.TurnHistory[0][KeysForgedThisTurn] = 1

	got := entriesAppended(g, func() { g.expireTurnScoped(0) })

	if len(got) != 0 {
		t.Errorf("rolling the turn history appended %v, want nothing", got)
	}
	h := g.State.TurnHistory[0]
	if h[KeysForgedLastTurn] != 1 || h[KeysForgedThisTurn] != 0 {
		t.Errorf("turn history after cleanup is last=%d this=%d, want 1 and 0",
			h[KeysForgedLastTurn], h[KeysForgedThisTurn])
	}
}
