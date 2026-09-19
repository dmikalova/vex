package engine

import "testing"

// LegalActions enumerates the root-action Commands available right now. These
// tests cover each phase and each kind of offer, and — most importantly — the
// contract that binds LegalActions to ApplyAction: every command it offers, the
// same state accepts (TestLegalActionsAreAllApplicable).

// countKind reports how many of cmds carry the given kind.
func countKind(cmds []Command, kind CommandKind) int {
	n := 0
	for _, c := range cmds {
		if c.Kind == kind {
			n++
		}
	}
	return n
}

// creatureUseCommandsFor returns the commands that use the creature id — the four
// kinds that name a creature in Card. Filtering by kind avoids matching a command
// whose Card is the zero LocalID for an unrelated reason (CommandEndTurn, or a
// hand play keyed by Hand), which would otherwise alias creature id 0.
func creatureUseCommandsFor(cmds []Command, id LocalID) []Command {
	var out []Command
	for _, c := range cmds {
		switch c.Kind {
		case CommandReap, CommandFight, CommandUseAction, CommandUnstun:
			if c.Card == id {
				out = append(out, c)
			}
		}
	}
	return out
}

func TestLegalActionsNoneForInactivePlayer(t *testing.T) {
	g := started(t) // player 0 is active
	if got := g.LegalActions(1); got != nil {
		t.Fatalf("inactive player has legal actions: %+v", got)
	}
}

func TestLegalActionsNoneWhenGameOver(t *testing.T) {
	g := started(t)
	g.State.Winner = 0
	if got := g.LegalActions(0); got != nil {
		t.Fatalf("a finished game has legal actions: %+v", got)
	}
}

func TestLegalActionsNoneInEngineDrivenPhase(t *testing.T) {
	g := started(t) // PhasePlay
	// The ready/draw/end phases resolve on their own and prompt through Requests,
	// not root actions, so no root action is legal in them.
	g.State.Phase = PhaseReady
	if got := g.LegalActions(0); got != nil {
		t.Fatalf("engine-driven phase has root actions: %+v", got)
	}
}

func TestLegalActionsHouseChoices(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.StartTurn(0)
	if g.State.Phase != PhaseChooseHouse {
		t.Fatalf("phase = %v, want PhaseChooseHouse", g.State.Phase)
	}
	cmds := g.LegalActions(0)
	// A fresh game leaves every house choosable, so every offer is a house choice
	// and Brobnar is among them.
	if len(cmds) == 0 {
		t.Fatal("no house choices offered")
	}
	sawBrobnar := false
	for _, c := range cmds {
		if c.Kind != CommandChooseHouse {
			t.Fatalf("non-house command in choose-house phase: %+v", c)
		}
		if c.House == Brobnar {
			sawBrobnar = true
		}
	}
	if !sawBrobnar {
		t.Fatal("Brobnar not offered among house choices")
	}
}

func TestLegalActionsHouseChoiceIsNoHouseWhenNoneAllowed(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	// Narrow the deck to one house, then bar it: the allowed set is empty, so the
	// only legal choice is No House.
	g.SetPlayerHouses(0, []House{Brobnar})
	g.StartTurn(0)
	g.State.HouseConstraints[0][0] = HouseConstraint{Kind: constraintCannotHouse, House: Brobnar}
	g.State.HouseConstraintCount[0] = 1
	cmds := g.LegalActions(0)
	want := []Command{{Kind: CommandChooseHouse, House: HouseNone}}
	if len(cmds) != 1 || cmds[0] != want[0] {
		t.Fatalf("no-house choice = %+v, want %+v", cmds, want)
	}
}

func TestLegalActionsEmptyBattlelineOffersOneFlank(t *testing.T) {
	g := started(t) // empty battleline
	id := g.AddToHand(testCreature("Lone", 3), 0)
	i := handIdxByID(g, 0, id)
	n := 0
	for _, c := range g.LegalActions(0) {
		if c.Kind == CommandPlayCreature && c.Hand == i {
			n++
		}
	}
	// On an empty line both flanks place the creature identically, so only one
	// play is offered rather than a redundant pair.
	if n != 1 {
		t.Fatalf("empty-battleline creature offered %d plays, want 1", n)
	}
}

func TestLegalActionsPopulatedBattlelineOffersBothFlanks(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("Anchor", 5), 0)
	id := g.AddToHand(testCreature("Newcomer", 3), 0)
	i := handIdxByID(g, 0, id)
	var flanks []bool
	for _, c := range g.LegalActions(0) {
		if c.Kind == CommandPlayCreature && c.Hand == i {
			flanks = append(flanks, c.Left)
		}
	}
	// With a creature already on the line the two flanks are distinct placements,
	// so both are offered.
	if len(flanks) != 2 || flanks[0] == flanks[1] {
		t.Fatalf("populated-battleline creature flanks = %v, want both left and right", flanks)
	}
}

func TestLegalActionsStunnedCreatureOffersOnlyUnstun(t *testing.T) {
	g := started(t)
	stunned := g.AddToBattleline(testCreature("Stunned", 3), 0)
	g.State.Cards[stunned].Stunned = true
	forStunned := creatureUseCommandsFor(g.LegalActions(0), stunned)
	if len(forStunned) != 1 || forStunned[0].Kind != CommandUnstun {
		t.Fatalf("stunned creature offers %+v, want a single Unstun", forStunned)
	}
}

func TestLegalActionsExhaustedStunnedCreatureOffersNothing(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("Spent", 3), 0)
	g.State.Cards[id].Stunned = true
	g.State.Cards[id].Exhausted = true
	// An exhausted creature cannot spend a use, so even a stunned one offers no
	// Unstun.
	if got := creatureUseCommandsFor(g.LegalActions(0), id); len(got) != 0 {
		t.Fatalf("exhausted stunned creature offered %+v, want nothing", got)
	}
}

// TestLegalActionsAreAllApplicable is the contract between the two halves of the
// root-action seam: every command LegalActions offers on a state, ApplyAction on
// that same state accepts. It builds a play phase rich enough to exercise every
// offer — each playable card type in hand, a creature to reap and fight, a
// creature with an Action, a stunned creature, an enemy to fight, and an artifact
// with an Action — then applies each offered command to an independent copy and
// fails on any the engine refuses.
func TestLegalActionsAreAllApplicable(t *testing.T) {
	g := started(t)
	// Hand: one of each playable type. The upgrade needs a host, added below.
	g.AddToHand(testCreature("HandCreature", 3), 0)
	g.AddToHand(exAutocannon(), 0) // artifact
	g.AddToHand(exBattleFury(), 0) // tactic
	g.AddToHand(exBruteStrength(), 0)
	// Battleline: a host for the upgrade, a plain user, a creature with an Action,
	// and a stunned creature.
	g.AddToBattleline(testCreature("Host", 5), 0)
	g.AddToBattleline(testCreature("Reaper", 3), 0)
	g.AddToBattleline(NewCard("Actor", Brobnar, Creature, Common,
		WithPower(3), WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 2})), 0)
	stunned := g.AddToBattleline(testCreature("Stunned", 3), 0)
	g.State.Cards[stunned].Stunned = true
	// An enemy creature gives the friendly creatures a legal fight target.
	g.AddToBattleline(testCreature("Enemy", 4), 1)
	// A friendly artifact with an Action.
	g.AddArtifact(NewCard("Relic", Brobnar, Artifact, Common,
		WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1})), 0)

	cmds := g.LegalActions(0)
	if len(cmds) == 0 {
		t.Fatal("expected a rich legal-action set, got none")
	}
	// The set must reach every offer kind, or the contract check is vacuous for it.
	for _, kind := range []CommandKind{
		CommandPlayCreature, CommandPlayArtifact, CommandPlayTactic, CommandPlayUpgrade,
		CommandDiscardFromHand, CommandReap, CommandFight, CommandUseAction,
		CommandUnstun, CommandEndTurn,
	} {
		if countKind(cmds, kind) == 0 {
			t.Errorf("legal set missing an offer of kind %d", kind)
		}
	}
	for _, cmd := range cmds {
		probe := *g
		probe.State = g.State.FastCopy()
		if err := probe.ApplyAction(cmd); err != nil {
			t.Errorf("LegalActions offered %+v but ApplyAction rejected it: %v", cmd, err)
		}
	}
}
