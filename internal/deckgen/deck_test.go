package deckgen

import (
	"strings"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// legalDeck builds a well-formed deck: three real-House pods, every slot holding a
// card of its pod's House.
func legalDeck() Deck {
	houses := [PodCount]engine.House{engine.Brobnar, engine.Dis, engine.Logos}
	var d Deck
	for i, h := range houses {
		d.Pods[i].House = h
		for j := range d.Pods[i].Slots {
			d.Pods[i].Slots[j] = Slot{Card: mkCard("C", h, engine.Common).Def}
		}
	}
	return d
}

// A well-formed deck validates without panicking. Degenerate-pool artifacts are
// passed over rather than rejected: a House-less (unfilled) pod, and an empty slot
// left by a pool too small to fill a real pod.
func TestValidateLegalDeck(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("validate panicked on a well-formed deck: %v", r)
		}
	}()
	legalDeck().validate()

	d := legalDeck()
	d.Pods[1] = HousePod{}      // House-less, empty slots — skipped, no panic.
	d.Pods[0].Slots[3] = Slot{} // empty slot in a real pod — skipped, no panic.
	d.validate()
}

// validate panics when a real pod holds a card of the wrong House.
func TestValidateMalformed(t *testing.T) {
	d := legalDeck()
	d.Pods[2].Slots[0] = Slot{Card: mkCard("X", engine.Untamed, engine.Common).Def}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected validate to panic on a mis-housed card")
		}
		if msg, _ := r.(string); !strings.Contains(msg, "not the pod's House") {
			t.Fatalf("panic = %v, want it to name the House mismatch", r)
		}
	}()
	d.validate()
}
