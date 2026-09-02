package engine

import (
	"strings"
	"testing"
)

// TestPlayRequirement covers the two Æmber gates a card can put on its own play:
// a threshold it only checks (Kelifi Dragon) and a cost it charges (Truebaru).
func TestPlayRequirement(t *testing.T) {
	if (PlayRequirement{}).text() != "" {
		t.Error("a card with no requirement should print no requirement line")
	}

	dragon := NewCard("Kelifi Dragon", Brobnar, Creature, Rare, WithPower(12),
		WithPlayRequirement(AemberThreshold(7)))
	wantDragon := "Kelifi Dragon cannot be played unless you have 7 Æmber or more."
	if text := RenderCardText(&dragon); !strings.Contains(text, wantDragon) {
		t.Errorf("threshold text = %q, want it to contain %q", text, wantDragon)
	}

	truebaru := NewCard("Truebaru", Dis, Creature, Rare, WithPower(7),
		WithPlayRequirement(AemberCost(3)))
	wantTruebaru := "You must lose 3 Æmber in order to play Truebaru."
	if text := RenderCardText(&truebaru); !strings.Contains(text, wantTruebaru) {
		t.Errorf("cost text = %q, want it to contain %q", text, wantTruebaru)
	}

	g := NewGame("A", "B", 1)
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	id := g.AddToHand(dragon, 0)

	g.SetAember(0, 6)
	if err := g.CanPlay(0, id); err != ErrPlayRequirement {
		t.Errorf("CanPlay under the threshold = %v, want %v", err, ErrPlayRequirement)
	}
	if _, err := g.PlayCreature(0, 0, false); err != ErrPlayRequirement {
		t.Errorf("PlayCreature under the threshold = %v, want %v", err, ErrPlayRequirement)
	}

	g.SetAember(0, 7)
	if err := g.CanPlay(0, id); err != nil {
		t.Errorf("CanPlay at the threshold = %v, want nil", err)
	}
	if _, err := g.PlayCreature(0, 0, false); err != nil {
		t.Fatalf("PlayCreature at the threshold = %v, want nil", err)
	}
	if g.Aember(0) != 7 {
		t.Errorf("a threshold should spend nothing, pool = %d, want 7", g.Aember(0))
	}

	// A cost is charged out of the pool as the card is played.
	h := NewGame("A", "B", 1)
	h.BeginTurn(0)
	if err := h.ChooseHouse(0, Dis); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	h.AddToHand(truebaru, 0)
	h.SetAember(0, 4)
	if _, err := h.PlayCreature(0, 0, false); err != nil {
		t.Fatalf("PlayCreature with the cost covered = %v, want nil", err)
	}
	if h.Aember(0) != 1 {
		t.Errorf("pool after paying the cost = %d, want 1", h.Aember(0))
	}
}
