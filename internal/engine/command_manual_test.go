package engine

import (
	"math/rand"
	"testing"
)

// ApplyManual performs each manual/debug force-edit by dispatching to the engine's
// Manual* method. One subtest per kind proves the dispatch lands on the right
// method and carries the command's fields through, so the whole manual vocabulary
// is covered from the one seam a driver applies a force-edit through.

func TestApplyManualSetManual(t *testing.T) {
	g := started(t)
	if err := g.ApplyManual(Command{
		Kind: CommandSetManual,
		Left: true,
	}, nil); err != nil {
		t.Fatalf("ApplyManual SetManual: %v", err)
	}
	if !g.Manual() {
		t.Fatal("SetManual command should turn manual mode on")
	}
	if err := g.ApplyManual(Command{
		Kind: CommandSetManual,
		Left: false,
	}, nil); err != nil {
		t.Fatalf("ApplyManual SetManual off: %v", err)
	}
	if g.Manual() {
		t.Fatal("SetManual command with Left=false should turn manual mode off")
	}
}

func TestApplyManualMove(t *testing.T) {
	g := started(t)
	id := g.AddToHand(testCreature("c", 3), 0)
	if err := g.ApplyManual(
		Command{
			Kind:  CommandManualMove,
			Card:  id,
			Index: int(ManualArchives),
		}, nil,
	); err != nil {
		t.Fatalf("ApplyManual Move: %v", err)
	}
	if got := g.Archives(0); len(got) != 1 || got[0] != id {
		t.Fatalf("archives = %v, want [%d]", got, id)
	}
}

func TestApplyManualReadyAndExhaust(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("c", 3), 0)
	if err := g.ApplyManual(Command{
		Kind: CommandManualExhaust,
		Card: id,
	}, nil); err != nil {
		t.Fatalf("ApplyManual Exhaust: %v", err)
	}
	if !g.Exhausted(id) {
		t.Fatal("ManualExhaust command should exhaust the card")
	}
	if err := g.ApplyManual(Command{
		Kind: CommandManualReady,
		Card: id,
	}, nil); err != nil {
		t.Fatalf("ApplyManual Ready: %v", err)
	}
	if g.Exhausted(id) {
		t.Fatal("ManualReady command should ready the card")
	}
}

func TestApplyManualAttach(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	buried := g.AddToHand(testCreature("buried", 2), 0)
	if err := g.ApplyManual(
		Command{
			Kind:  CommandManualAttach,
			Card:  host,
			Card2: buried,
			Left:  true,
		}, nil,
	); err != nil {
		t.Fatalf("ApplyManual Attach: %v", err)
	}
	if u, ok := g.firstUnder(host); !ok || u != buried {
		t.Fatalf("firstUnder = %d,%v, want %d", u, ok, buried)
	}
	if !g.State.Cards[buried].UnderFaceDown {
		t.Error("Left=true should place the card face down")
	}
}

func TestApplyManualPlace(t *testing.T) {
	g := started(t)
	left := g.AddToBattleline(testCreature("left", 3), 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	mid := g.AddToHand(testCreature("mid", 3), 0)
	if err := g.ApplyManual(
		Command{
			Kind:  CommandManualPlace,
			Card:  mid,
			Index: 1,
		}, nil,
	); err != nil {
		t.Fatalf("ApplyManual Place: %v", err)
	}
	if line := g.Battleline(
		0,
	); len(line) != 3 || line[0] != left || line[1] != mid ||
		line[2] != right {
		t.Fatalf("battleline = %v, want [left mid right]", line)
	}
}

func TestApplyManualDetach(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	up := g.Register(exBruteStrength(), 0)
	g.AttachUpgrade(host, up)
	if err := g.ApplyManual(Command{
		Kind: CommandManualDetach,
		Card: up,
	}, nil); err != nil {
		t.Fatalf("ApplyManual Detach: %v", err)
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != up {
		t.Fatalf("hand = %v, want [%d]", g.Hand(0), up)
	}
}

func TestApplyManualAmber(t *testing.T) {
	g := started(t)
	if err := g.ApplyManual(
		Command{
			Kind:   CommandManualAmber,
			Player: 1,
			Delta:  4,
		}, nil,
	); err != nil {
		t.Fatalf("ApplyManual Amber: %v", err)
	}
	if g.Aember(1) != 4 {
		t.Fatalf("Æmber = %d, want 4", g.Aember(1))
	}
}

func TestApplyManualChains(t *testing.T) {
	g := started(t)
	if err := g.ApplyManual(
		Command{
			Kind:   CommandManualChains,
			Player: 0,
			Delta:  3,
		}, nil,
	); err != nil {
		t.Fatalf("ApplyManual Chains: %v", err)
	}
	if g.State.Chains[0] != 3 {
		t.Fatalf("chains = %d, want 3", g.State.Chains[0])
	}
}

func TestApplyManualForgeAndUnforge(t *testing.T) {
	g := started(t)
	if err := g.ApplyManual(
		Command{
			Kind:   CommandManualForgeColor,
			Player: 0,
			Index:  int(KeyColorRed),
		}, nil,
	); err != nil {
		t.Fatalf("ApplyManual ForgeColor: %v", err)
	}
	if g.Keys(0) != 1 || g.KeyColors(0)[0] != KeyColorRed {
		t.Fatalf("keys = %d colors = %v, want 1 red", g.Keys(0), g.KeyColors(0))
	}
	if err := g.ApplyManual(Command{
		Kind:   CommandManualUnforge,
		Player: 0,
	}, nil); err != nil {
		t.Fatalf("ApplyManual Unforge: %v", err)
	}
	if g.Keys(0) != 0 {
		t.Fatalf("keys = %d, want 0 after unforge", g.Keys(0))
	}
}

func TestApplyManualHouse(t *testing.T) {
	g := started(t)
	if err := g.ApplyManual(Command{
		Kind:  CommandManualHouse,
		House: Dis,
	}, nil); err != nil {
		t.Fatalf("ApplyManual House: %v", err)
	}
	if g.State.ActiveHouse != Dis {
		t.Fatalf("active house = %v, want Dis", g.State.ActiveHouse)
	}
}

func TestApplyManualAddCard(t *testing.T) {
	g := started(t)
	resolve := func(name string) (CardDefinition, bool) {
		if name != "Import" {
			return CardDefinition{}, false
		}
		return NewCard("Import", Logos, Tactic, Common), true
	}
	before := len(g.Hand(1))
	if err := g.ApplyManual(
		Command{
			Kind:   CommandManualAddCard,
			Name:   "Import",
			Player: 1,
		}, resolve,
	); err != nil {
		t.Fatalf("ApplyManual AddCard: %v", err)
	}
	if len(g.Hand(1)) != before+1 {
		t.Fatalf("hand size = %d, want %d", len(g.Hand(1)), before+1)
	}

	// A name the pool does not know fails loudly rather than adding nothing quietly,
	// so a diverged replay is caught.
	err := g.ApplyManual(Command{
		Kind: CommandManualAddCard,
		Name: "Unknown",
	}, resolve)
	if err == nil {
		t.Fatal("ApplyManual AddCard with an unknown name should error")
	}
}

func TestApplyManualRejectsNonManual(t *testing.T) {
	g := started(t)
	// A root action or a choice answer is not a manual edit: ApplyManual refuses it
	// so the two dispatch seams stay disjoint.
	if err := g.ApplyManual(Command{Kind: CommandReap}, nil); err != ErrNotManualAction {
		t.Fatalf("ApplyManual on a root action = %v, want ErrNotManualAction", err)
	}
}

// TestLegalActionsNeverOffersManualKinds pins the sim-safety invariant behind
// splitting the manual force-edits out as their own Command kinds: LegalActions
// enumerates only real root actions, never a manual edit, so the canonical turn
// loop (RunMatch) and the simulator that drives it can never wander into manual
// mode. It plays many seeded games with a random legal-action driver and asserts
// every command offered at every step is a root action, not a manual kind.
func TestLegalActionsNeverOffersManualKinds(t *testing.T) {
	for seed := int64(1); seed <= 30; seed++ {
		g := NewGame("Alice", "Bob", seed)
		g.StartGame(0)
		rng := rand.New(rand.NewSource(seed))
		for step := 0; g.Winner() < 0 && step < 400; step++ {
			if g.State.Phase == PhaseEndOfTurn {
				g.StartTurn(1 - g.State.ActivePlayer)
				continue
			}
			p := g.State.ActivePlayer
			cmds := g.LegalActions(p)
			if len(cmds) == 0 {
				t.Fatalf("seed %d step %d: no legal actions in phase %v", seed, step, g.State.Phase)
			}
			for _, c := range cmds {
				if c.Kind >= CommandSetManual {
					t.Fatalf("seed %d: LegalActions offered manual kind %d — the sim would "+
						"drive manual mode", seed, c.Kind)
				}
			}
			_ = g.ApplyAction(cmds[rng.Intn(len(cmds))])
		}
	}
}
