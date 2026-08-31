package engine

import (
	"strings"
	"testing"
)

// A sound state — creatures in play, an attached upgrade, and cards in hand —
// reports no invariant violation. This also exercises the card-conservation walk
// over both players' zones and a creature's upgrades.
func TestInvariantErrorSound(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("Guard", 5), 0)
	attachUpgrade(g, c, exBruteStrength())
	g.AddToHand(testCreature("Spare", 3), 1)
	if err := g.InvariantError(); err != nil {
		t.Fatalf("sound state reported an invariant error: %v", err)
	}
}

func TestInvariantErrorEconomy(t *testing.T) {
	cases := []struct {
		name string
		mut  func(g *Game)
		want string
	}{
		{"negative aember", func(g *Game) { g.State.Aember[0] = -1 }, "negative Æmber"},
		{
			"too many keys",
			func(g *Game) { g.State.Keys[1] = KeysToWin + 1 },
			"out-of-range key count",
		},
		{"negative chains", func(g *Game) { g.State.Chains[0] = -1 }, "negative chains"},
		{"winner out of range", func(g *Game) { g.State.Winner = 2 }, "winner is out of range"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGame("A", "B", 1)
			tc.mut(g)
			err := g.InvariantError()
			if err == nil {
				t.Fatalf("expected an invariant error, got nil")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not mention %q", err, tc.want)
			}
		})
	}
}

// A card sitting in two zones at once breaks conservation.
func TestInvariantErrorCardDuplicated(t *testing.T) {
	g := NewGame("A", "B", 1)
	id := g.AddToHand(testCreature("Ghost", 4), 0)
	g.State.Discard[0].add(id) // now in hand and discard both
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a duplicated card, got nil")
	}
	if !strings.Contains(err.Error(), "want exactly 1") {
		t.Fatalf("error %q does not describe a conservation break", err)
	}
}

// A card that carries a host back-link but is not held in any creature's chain is
// a dangling upgrade — it is "in play" as an attachment with nothing to attach to.
func TestInvariantErrorDanglingUpgrade(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("Guard", 5), 0)
	up := g.AddToHand(exBruteStrength(), 0)
	g.State.Cards[up].HostPlus = upgradePlus(host) // back-link with no chain entry
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a dangling upgrade, got nil")
	}
	if !strings.Contains(err.Error(), "dangling upgrade") {
		t.Fatalf("error %q does not describe a dangling upgrade", err)
	}
}

// A non-Upgrade card threaded into a creature's chain violates the chain's type.
func TestInvariantErrorNonUpgradeInChain(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("Guard", 5), 0)
	intruder := g.Register(testCreature("Sneak", 2), 0)
	g.AttachUpgrade(host, intruder)
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a non-Upgrade in a chain, got nil")
	}
	if !strings.Contains(err.Error(), "not an Upgrade") {
		t.Fatalf("error %q does not describe a wrong-type chain entry", err)
	}
}

// An upgrade in a host's chain whose back-link points somewhere else is a broken
// attachment even though conservation still counts it once.
func TestInvariantErrorUpgradeBackLinkDisagrees(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("Guard", 5), 0)
	up := g.Register(exBruteStrength(), 0)
	g.AttachUpgrade(host, up)
	g.State.Cards[up].HostPlus = upgradePlus(host) + 1 // point the back-link elsewhere
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a disagreeing back-link, got nil")
	}
	if !strings.Contains(err.Error(), "back-link disagrees") {
		t.Fatalf("error %q does not describe a back-link mismatch", err)
	}
}
