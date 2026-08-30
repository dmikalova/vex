package engine

import "testing"

// AddToDiscard and AddToArchives register a card and place it in the right pile.
func TestAddToDiscardAndArchives(t *testing.T) {
	g := NewGame("A", "B", 1)
	d := g.AddToDiscard(testCreature("Buried", 3), 0)
	if !g.State.Discard[0].contains(d) {
		t.Errorf("AddToDiscard did not place the card in the discard pile")
	}
	a := g.AddToArchives(testCreature("Stashed", 4), 1)
	if !g.State.Archives[1].contains(a) {
		t.Errorf("AddToArchives did not place the card in the archives")
	}
}
