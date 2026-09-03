package deckgen

import (
	"math/rand"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

func connCard(name string, h engine.House, conn Connection) Card {
	c := mkCard(name, h, engine.Common)
	c.Profile.Connection = conn
	return c
}

func connectedCard(name string, h engine.House) Card {
	return mkCard(name, h, engine.Connected)
}

// pulls is the plain connection: one copy of each named card, every time.
func pulls(names ...string) Connection {
	conn := Connection{}
	for _, n := range names {
		conn.Cards = append(conn.Cards, ConnectedCard{Name: n, Copies: 1, Chance: 1})
	}
	return conn
}

func gen(set Set) *generator {
	return &generator{set: set, r: rand.New(rand.NewSource(1)), placed: map[string]bool{}}
}

func countName(pod HousePod, name string) int {
	n := 0
	for _, s := range pod.Slots {
		if s.Card.Name == name {
			n++
		}
	}
	return n
}

// A puller placed in a pod pulls exactly one copy of its connected partner.
func TestConnectionPullsPartner(t *testing.T) {
	set := NewSet("S", []Card{
		connCard("Puller", engine.Brobnar, pulls("Partner")),
		mkCard("Filler", engine.Brobnar, engine.Common),
		connectedCard("Partner", engine.Brobnar),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	deck := Generate(set, 1)
	total := 0
	for _, pod := range deck.Pods {
		total += countName(pod, "Partner")
	}
	if total != 1 {
		t.Fatalf("Partner count = %d, want exactly 1", total)
	}
}

// A Connected card never rolls into a deck without a puller.
func TestConnectedExcludedFromPool(t *testing.T) {
	set := NewSet("S", []Card{
		mkCard("Filler", engine.Brobnar, engine.Common),
		connectedCard("Lonely", engine.Brobnar),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	deck := Generate(set, 3)
	for _, pod := range deck.Pods {
		if countName(pod, "Lonely") != 0 {
			t.Fatal("a Connected card rolled without a puller")
		}
	}
}

// A maverick puller does not fire its connection (cross-house is off for v1).
func TestConnectionMaverickPullerSkipped(t *testing.T) {
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, pulls("Q")),
		connectedCard("Q", engine.Brobnar),
	}, DefaultTuning())
	g := gen(set)
	pod := HousePod{House: engine.Brobnar}
	pod.Slots[0] = Slot{Card: set.byName["P"].Def, Maverick: true}
	out := g.expandConnections(pod)
	if countName(out, "Q") != 0 {
		t.Fatal("maverick puller should not pull a partner")
	}
}

// A connection naming a card absent from the set is an authoring error: NewSet
// fails loudly rather than silently dropping the link.
func TestConnectionMissingPartner(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewSet must panic on a connection to a missing card")
		}
	}()
	NewSet("S", []Card{
		connCard("P", engine.Brobnar, pulls("Ghost")),
	}, DefaultTuning())
}

// A partner already in the pod is left as the single copy, not duplicated.
func TestConnectionPartnerAlreadyPresent(t *testing.T) {
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, pulls("Q")),
		connectedCard("Q", engine.Brobnar),
	}, DefaultTuning())
	g := gen(set)
	pod := HousePod{House: engine.Brobnar}
	pod.Slots[0] = Slot{Card: set.byName["P"].Def}
	pod.Slots[1] = Slot{Card: set.byName["Q"].Def}
	out := g.expandConnections(pod)
	if countName(out, "Q") != 1 {
		t.Fatalf("Q count = %d, want 1 (already present, not duplicated)", countName(out, "Q"))
	}
}

// When every slot is a protected puller there is no room to place a partner.
func TestConnectionNoFreeSlot(t *testing.T) {
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, pulls("Q")),
		connectedCard("Q", engine.Brobnar),
	}, DefaultTuning())
	g := gen(set)
	pod := HousePod{House: engine.Brobnar}
	for i := range pod.Slots {
		pod.Slots[i] = Slot{Card: set.byName["P"].Def}
	}
	out := g.expandConnections(pod)
	if countName(out, "Q") != 0 {
		t.Fatal("a full pod of pullers has no room for a partner")
	}
}

// A pulled partner's own connection resolves on a later pass (fixpoint).
func TestConnectionFixpointChain(t *testing.T) {
	q := connectedCard("Q", engine.Brobnar)
	q.Profile.Connection = pulls("R")
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, pulls("Q")),
		q,
		connectedCard("R", engine.Brobnar),
	}, DefaultTuning())
	g := gen(set)
	pod := HousePod{House: engine.Brobnar}
	pod.Slots[0] = Slot{Card: set.byName["P"].Def}
	out := g.expandConnections(pod)
	if countName(out, "Q") != 1 || countName(out, "R") != 1 {
		t.Fatalf("chain: Q=%d R=%d, want 1 and 1", countName(out, "Q"), countName(out, "R"))
	}
}

// A connection asking for several copies places every one of them.
func TestConnectionPullsSeveralCopies(t *testing.T) {
	conn := Connection{Cards: []ConnectedCard{{Name: "Q", Copies: 3, Chance: 1}}}
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, conn),
		connectedCard("Q", engine.Brobnar),
	}, DefaultTuning())
	g := gen(set)
	pod := HousePod{House: engine.Brobnar}
	pod.Slots[0] = Slot{Card: set.byName["P"].Def}
	out := g.expandConnections(pod)
	if countName(out, "Q") != 3 {
		t.Fatalf("Q count = %d, want 3", countName(out, "Q"))
	}
}

// A copy already in the pod counts toward the connection's total, so a pod that
// rolled one of the pulled card is topped up rather than filled from scratch.
func TestConnectionCopiesCountExisting(t *testing.T) {
	conn := Connection{Cards: []ConnectedCard{{Name: "Q", Copies: 2, Chance: 1}}}
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, conn),
		connectedCard("Q", engine.Brobnar),
	}, DefaultTuning())
	g := gen(set)
	pod := HousePod{House: engine.Brobnar}
	pod.Slots[0] = Slot{Card: set.byName["P"].Def}
	pod.Slots[1] = Slot{Card: set.byName["Q"].Def}
	out := g.expandConnections(pod)
	if countName(out, "Q") != 2 {
		t.Fatalf("Q count = %d, want 2", countName(out, "Q"))
	}
}

// A chancy pull fires only sometimes: over many pods it appears in some and not
// in others, and it is never pulled more than its copy count.
func TestConnectionChancePullIsSometimes(t *testing.T) {
	conn := Connection{Cards: []ConnectedCard{{Name: "Q", Copies: 1, Chance: 0.5}}}
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, conn),
		connectedCard("Q", engine.Brobnar),
	}, DefaultTuning())
	g := gen(set)
	hits := 0
	for i := 0; i < 200; i++ {
		pod := HousePod{House: engine.Brobnar}
		pod.Slots[0] = Slot{Card: set.byName["P"].Def}
		n := countName(g.expandConnections(pod), "Q")
		if n > 1 {
			t.Fatalf("Q count = %d, want at most 1", n)
		}
		hits += n
	}
	if hits == 0 || hits == 200 {
		t.Fatalf("chancy pull fired %d/200 times, want sometimes", hits)
	}
}

// A connection pulling fewer than one copy, or at an impossible rate, is an
// authoring error NewSet reports rather than silently normalizing.
func TestConnectionInvalidPullPanics(t *testing.T) {
	for name, conn := range map[string]Connection{
		"no copies":    {Cards: []ConnectedCard{{Name: "Q", Copies: 0, Chance: 1}}},
		"never fires":  {Cards: []ConnectedCard{{Name: "Q", Copies: 1, Chance: 0}}},
		"always twice": {Cards: []ConnectedCard{{Name: "Q", Copies: 1, Chance: 2}}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("NewSet must panic on an impossible connection")
				}
			}()
			NewSet("S", []Card{
				connCard("P", engine.Brobnar, conn),
				connectedCard("Q", engine.Brobnar),
			}, DefaultTuning())
		})
	}
}
