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
		connCard("Puller", engine.Brobnar, Connection{Cards: []string{"Partner"}}),
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
		connCard("P", engine.Brobnar, Connection{Cards: []string{"Q"}}),
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
		connCard("P", engine.Brobnar, Connection{Cards: []string{"Ghost"}}),
	}, DefaultTuning())
}

// A partner already in the pod is left as the single copy, not duplicated.
func TestConnectionPartnerAlreadyPresent(t *testing.T) {
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, Connection{Cards: []string{"Q"}}),
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
		connCard("P", engine.Brobnar, Connection{Cards: []string{"Q"}}),
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
	q.Profile.Connection = Connection{Cards: []string{"R"}}
	set := NewSet("S", []Card{
		connCard("P", engine.Brobnar, Connection{Cards: []string{"Q"}}),
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
