package engine

import (
	"errors"
	"testing"
)

// ApplyAction performs each root action by dispatching to the engine method the
// matching player click calls. One subtest per kind proves the dispatch lands on
// the right method and reports its result, so the whole root-action vocabulary is
// covered from the one seam a driver uses.

func TestApplyActionChooseHouse(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.StartTurn(0)
	if g.State.Phase != PhaseChooseHouse {
		t.Fatalf("phase after StartTurn = %v, want PhaseChooseHouse", g.State.Phase)
	}
	if err := g.ApplyAction(Command{
		Kind:  CommandChooseHouse,
		House: Brobnar,
	}); err != nil {
		t.Fatalf("ApplyAction ChooseHouse: %v", err)
	}
	if g.State.ActiveHouse != Brobnar {
		t.Fatalf("active house = %v, want Brobnar", g.State.ActiveHouse)
	}
}

func TestApplyActionPlayCreature(t *testing.T) {
	g := started(t)
	id := g.AddToHand(testCreature("Ganger", 4), 0)
	if err := g.ApplyAction(Command{
		Kind: CommandPlayCreature, Hand: handIdxByID(g, 0, id), Left: true,
	}); err != nil {
		t.Fatalf("ApplyAction PlayCreature: %v", err)
	}
	if got := g.Battleline(0); len(got) != 1 || got[0] != id {
		t.Fatalf("battleline = %v, want [%d]", got, id)
	}
}

func TestApplyActionPlayArtifact(t *testing.T) {
	g := started(t)
	id := g.AddToHand(exAutocannon(), 0)
	if err := g.ApplyAction(Command{
		Kind: CommandPlayArtifact, Hand: handIdxByID(g, 0, id),
	}); err != nil {
		t.Fatalf("ApplyAction PlayArtifact: %v", err)
	}
	if got := g.Artifacts(0); len(got) != 1 || got[0] != id {
		t.Fatalf("artifact row = %v, want [%d]", got, id)
	}
}

func TestApplyActionPlayTactic(t *testing.T) {
	g := started(t)
	id := g.AddToHand(exBattleFury(), 0)
	before := g.Aember(0)
	if err := g.ApplyAction(Command{
		Kind: CommandPlayTactic, Hand: handIdxByID(g, 0, id),
	}); err != nil {
		t.Fatalf("ApplyAction PlayTactic: %v", err)
	}
	// Battle Fury carries an Æmber bonus, so playing it gains a pip.
	if g.Aember(0) != before+1 {
		t.Fatalf("Æmber after Tactic = %d, want %d", g.Aember(0), before+1)
	}
}

func TestApplyActionPlayUpgrade(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("Host", 5), 0)
	up := g.AddToHand(exBruteStrength(), 0)
	if err := g.ApplyAction(Command{
		Kind: CommandPlayUpgrade, Hand: handIdxByID(g, 0, up),
	}); err != nil {
		t.Fatalf("ApplyAction PlayUpgrade: %v", err)
	}
	if got := g.Upgrades(host); len(got) != 1 || got[0] != up {
		t.Fatalf("upgrades on host = %v, want [%d]", got, up)
	}
}

func TestApplyActionDiscardFromHand(t *testing.T) {
	g := started(t)
	id := g.AddToHand(testCreature("Chaff", 2), 0)
	if err := g.ApplyAction(Command{
		Kind: CommandDiscardFromHand, Hand: handIdxByID(g, 0, id),
	}); err != nil {
		t.Fatalf("ApplyAction DiscardFromHand: %v", err)
	}
	if got := g.State.Discard[0].slice(); len(got) != 1 || got[0] != id {
		t.Fatalf("discard = %v, want [%d]", got, id)
	}
}

func TestApplyActionReap(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("Reaper", 3), 0)
	before := g.Aember(0)
	if err := g.ApplyAction(Command{
		Kind: CommandReap,
		Card: g.Battleline(0)[0],
	}); err != nil {
		t.Fatalf("ApplyAction Reap: %v", err)
	}
	if g.Aember(0) != before+1 {
		t.Fatalf("Æmber after reap = %d, want %d", g.Aember(0), before+1)
	}
}

func TestApplyActionUnstun(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("Stunned", 3), 0)
	g.State.Cards[id].Stunned = true
	if err := g.ApplyAction(Command{
		Kind: CommandUnstun,
		Card: id,
	}); err != nil {
		t.Fatalf("ApplyAction Unstun: %v", err)
	}
	if g.State.Cards[id].Stunned {
		t.Fatal("creature is still stunned after Unstun")
	}
}

func TestApplyActionUseAction(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(NewCard("Actor", Brobnar, Creature, Common,
		WithPower(3), WithAbility(TriggerAction, GainAember{
			Player: Controller,
			Amount: 5,
		})), 0)
	if err := g.ApplyAction(Command{
		Kind: CommandUseAction,
		Card: id,
	}); err != nil {
		t.Fatalf("ApplyAction UseAction: %v", err)
	}
	if g.Aember(0) != 5 {
		t.Fatalf("Æmber after Action = %d, want 5", g.Aember(0))
	}
}

func TestApplyActionFight(t *testing.T) {
	g := started(t)
	attacker := g.AddToBattleline(testCreature("Attacker", 3), 0)
	defender := g.AddToBattleline(testCreature("Defender", 5), 1)
	if err := g.ApplyAction(Command{
		Kind: CommandFight, Card: attacker, Card2: defender,
	}); err != nil {
		t.Fatalf("ApplyAction Fight: %v", err)
	}
	if g.Damage(defender) != 3 {
		t.Fatalf("defender damage = %d, want 3", g.Damage(defender))
	}
}

func TestApplyActionEndTurn(t *testing.T) {
	g := started(t)
	if err := g.ApplyAction(Command{Kind: CommandEndTurn}); err != nil {
		t.Fatalf("ApplyAction EndTurn: %v", err)
	}
	if g.State.ActivePlayer != 1 {
		t.Fatalf("active player after end turn = %d, want 1", g.State.ActivePlayer)
	}
}

// An illegal play is not swallowed: ApplyAction returns the engine's own error, so
// a driver can surface why the click was refused.
func TestApplyActionSurfacesEngineError(t *testing.T) {
	g := started(t)
	offHouse := g.AddToHand(NewCard("Untamed One", Untamed, Creature, Common, WithPower(2)), 0)
	err := g.ApplyAction(Command{
		Kind: CommandPlayCreature, Hand: handIdxByID(g, 0, offHouse), Left: true,
	})
	if !errors.Is(err, ErrWrongHouse) {
		t.Fatalf("playing an off-house creature returned %v, want ErrWrongHouse", err)
	}
}

// A choice-answer Command is not a root action; ApplyAction rejects it rather than
// silently doing nothing.
func TestApplyActionRejectsNonRoot(t *testing.T) {
	g := started(t)
	if err := g.ApplyAction(Command{
		Kind: CommandPickCard,
		Card: 1,
	}); !errors.Is(
		err, ErrNotRootAction,
	) {
		t.Fatalf("ApplyAction on a card-pick command returned %v, want ErrNotRootAction", err)
	}
}

// omegaDriver is a scripted ActionChooser for the turn-loop test: it plays an
// Omega creature when it can (to end a play phase mid-action), names Brobnar when
// asked for a house, and otherwise takes the first legal action.
type omegaDriver struct{}

func (omegaDriver) ChooseCreature(_, _ string, c []LocalID) (LocalID, bool) {
	if len(c) == 0 {
		return 0, false
	}
	return c[0], true
}

func (omegaDriver) ChooseAction(actions []Command) Command {
	for _, a := range actions {
		if a.Kind == CommandPlayCreature {
			return a
		}
	}
	for _, a := range actions {
		if a.Kind == CommandChooseHouse && a.House == Brobnar {
			return a
		}
	}
	return actions[0]
}

// RunMatch deals the game, drives it through the ActionChooser, and hands off a
// turn an effect ended without a paired StartTurn. Here player 0 plays an Omega
// creature — ending their play phase mid-action — and the loop's hand-off starts
// player 1's turn, whose forge phase forges their third key and wins. This walks
// every branch of the loop: setup, a house choice, a play, the Omega hand-off, and
// the win that stops it.
func TestRunMatchDrivesToWinnerThroughEndedTurn(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	omega := testCreature("Omega", 3, WithKeywords(Omega))
	for range 20 {
		g.AddToDeck(omega, 0)
		g.AddToDeck(omega, 1)
	}
	// Player 1 forges their winning key the moment their turn begins.
	g.State.ForgeCanonicalKeys(1, KeysToWin-1)
	g.State.Aember[1] = KeyCost
	g.SetChooser(0, omegaDriver{})
	g.SetChooser(1, omegaDriver{})

	g.RunMatch(0)

	if g.Winner() != 1 {
		t.Fatalf("winner = %d, want 1", g.Winner())
	}
}

func TestRunMatchSkipsPhaseEndOfTurnHandOff(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.StartGame(0)
	g.State.Phase = PhaseEndOfTurn
	g.State.ActivePlayer = 0
	g.SetChooser(0, omegaDriver{})
	g.SetChooser(1, omegaDriver{})

	g.RunMatch(0)

	if g.State.ActivePlayer != 1 {
		t.Fatalf("active player after end-of-turn handoff = %d, want 1", g.State.ActivePlayer)
	}
}
