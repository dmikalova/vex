package engine

import (
	"slices"
	"strings"
	"testing"
)

// DeckList reads the whole catalog, so it lists a player's entire deck — draw
// pile included — deduplicated with an "xN" count and sorted.
func TestDeckList(t *testing.T) {
	g := started(t)
	// Player 0 owns two copies of Ape (an "xN" entry) and one Ogre (a single).
	g.AddToDeck(testCreature("Ape", 4), 0)
	g.AddToDeck(testCreature("Ape", 4), 0)
	g.AddToDeck(testCreature("Ogre", 6), 0)
	// Player 1 owns a single Imp, which must not appear in player 0's list.
	g.AddToDeck(testCreature("Imp", 1), 1)

	list := g.DeckList(0)
	want := []string{"Ape x2", "Ogre"}
	if !slices.Equal(list, want) {
		t.Errorf("DeckList(0) = %v, want %v", list, want)
	}
	if !slices.IsSorted(list) {
		t.Errorf("DeckList(0) is not sorted: %v", list)
	}
	if strings.Join(g.DeckList(1), ",") != "Imp" {
		t.Errorf("DeckList(1) = %v, want [Imp]", g.DeckList(1))
	}
}
