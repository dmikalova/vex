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
	// The normal-build no-op runs at every turn boundary; call it directly so the
	// empty body is exercised without an -tags assert build.
	g.assertInvariants()
}

func TestInvariantErrorEconomy(t *testing.T) {
	cases := []struct {
		name string
		mut  func(g *Game)
		want string
	}{
		{"negative aember", func(g *Game) { g.State.Aember[0] = -1 }, "negative Æmber"},
		{
			"key forged after a gap",
			func(g *Game) { g.State.KeyColors[1][MaxKeys-1] = KeyColorRed },
			"forged after an unforged slot",
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

// Æmber on a card is a saturating counter: it may never go negative, which would
// otherwise leak into a pool when the card leaves play.
func TestInvariantErrorNegativeAmberOnCard(t *testing.T) {
	g := NewGame("A", "B", 1)
	id := g.AddToHand(testCreature("Ghost", 4), 0)
	g.State.Cards[id].Amber = -1
	err := g.InvariantError()
	if err == nil {
		t.Fatal("expected an invariant error for negative Æmber on a card, got nil")
	}
	if !strings.Contains(err.Error(), "negative Æmber on it") {
		t.Fatalf("error %q does not mention negative Æmber on a card", err)
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

// A non-Upgrade, non-Creature card threaded into a creature's chain violates the
// chain's type. A creature is a valid chain entry (played as an upgrade, ADR 0026),
// so the intruder here is a Tactic, which can never occupy a chain.
func TestInvariantErrorNonUpgradeInChain(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("Guard", 5), 0)
	intruder := g.Register(
		NewCard("Sneak", Brobnar, Tactic, Common), 0)
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

// A card that carries an under-host back-link but is not held in any host's
// under-chain is a dangling under-card — mirroring the dangling-upgrade case
// above for the Under mechanic (see ADR 0016).
func TestInvariantErrorDanglingUnderCard(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	buried := g.AddToHand(testCreature("Buried", 2), 0)
	g.State.Cards[buried].UnderHostPlus = underPlus(host) // back-link with no chain entry
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a dangling under-card, got nil")
	}
	if !strings.Contains(err.Error(), "dangling under-card") {
		t.Fatalf("error %q does not describe a dangling under-card", err)
	}
}

// A card in a host's under-chain whose back-link points somewhere else is a
// broken attachment even though conservation still counts it once.
func TestInvariantErrorUnderBackLinkDisagrees(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	other := g.AddArtifact(NewCard("Other", Brobnar, Artifact, Common), 0)
	buried := g.Register(testCreature("Buried", 2), 0)
	g.AttachUnder(host, buried, false)
	g.State.Cards[buried].UnderHostPlus = underPlus(other) // point the back-link elsewhere
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a disagreeing under back-link, got nil")
	}
	if !strings.Contains(err.Error(), "back-link disagrees") {
		t.Fatalf("error %q does not describe a back-link mismatch", err)
	}
}

// A creature whose damage has caught up with its power has no business sitting in
// the battleline: some state-based sweep failed to notice it.
func TestInvariantErrorDamagedCreatureStillInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	id := g.AddToBattleline(testCreature("Oak", 3), 0)
	g.State.Cards[id].Damage = 3
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a creature at or below its damage, got nil")
	}
	if !strings.Contains(err.Error(), "want power above damage") {
		t.Fatalf("error %q does not describe damage catching up with power", err)
	}
}

// Per-match state belongs to cards in play, so a card in hand still carrying
// damage means a leave-play path forgot to reset it.
func TestInvariantErrorOutOfPlayCardKeepsItsState(t *testing.T) {
	g := NewGame("A", "B", 1)
	id := g.AddToHand(testCreature("Oak", 3), 0)
	g.State.Cards[id].Damage = 1
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for stranded in-play state, got nil")
	}
	if !strings.Contains(err.Error(), "still carries in-play state") {
		t.Fatalf("error %q does not describe stranded in-play state", err)
	}
}

// giganticBase and giganticArt build the two halves of a gigantic named for a
// test: the base carries the stats and holds the battleline slot, the art carries
// only its role.
func giganticBase(name string, power int) CardDefinition {
	def := NewCard(name, Sanctum, Creature, Rare, WithPower(power))
	def.GiganticRole = GiganticBase
	return def
}

func giganticArt(name string) CardDefinition {
	def := NewCard(name, Sanctum, Creature, Rare)
	def.GiganticRole = GiganticArt
	return def
}

// A linked gigantic — base in the battleline, slot-less art half linked to it —
// is sound: conservation counts the art half through its base (ADR 0042).
func TestInvariantErrorLinkedGiganticIsSound(t *testing.T) {
	g := NewGame("A", "B", 1)
	base := g.AddToBattleline(giganticBase("Deusillus", 15), 0)
	art := g.Register(giganticArt("Deusillus"), 0)
	g.linkGiganticPartners(base, art)
	if err := g.InvariantError(); err != nil {
		t.Fatalf("linked gigantic reported an invariant error: %v", err)
	}
}

// An art half carrying a partner link with no base holding it is a dangling
// gigantic — the leave-play funnel must clear the link on both halves together.
func TestInvariantErrorDanglingGigantic(t *testing.T) {
	g := NewGame("A", "B", 1)
	base := g.AddToBattleline(giganticBase("Deusillus", 15), 0)
	art := g.AddToHand(giganticArt("Deusillus"), 0)
	g.State.Cards[art].GiganticPartnerPlus = giganticPlus(base) // link with no base holding it
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a dangling gigantic, got nil")
	}
	if !strings.Contains(err.Error(), "dangling gigantic") {
		t.Fatalf("error %q does not describe a dangling gigantic", err)
	}
}

// A linked gigantic whose art half's back-link points elsewhere is a broken pair
// even though conservation still counts it once.
func TestInvariantErrorGiganticBackLinkDisagrees(t *testing.T) {
	g := NewGame("A", "B", 1)
	base := g.AddToBattleline(giganticBase("Deusillus", 15), 0)
	art := g.Register(giganticArt("Deusillus"), 0)
	g.linkGiganticPartners(base, art)
	g.State.Cards[art].GiganticPartnerPlus = giganticPlus(base) + 1 // point the back-link elsewhere
	err := g.InvariantError()
	if err == nil {
		t.Fatalf("expected an invariant error for a disagreeing gigantic back-link, got nil")
	}
	if !strings.Contains(err.Error(), "partner back-link disagrees") {
		t.Fatalf("error %q does not describe a back-link mismatch", err)
	}
}
