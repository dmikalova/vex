package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Orator Hissaro
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Dinosaur • Politician
//
//	Deploy.
//	Play: Ready and exalt each neighboring Creature. For the remainder of the turn, those Creatures belong to house Saurian.
func TestOratorHissaro(t *testing.T) {
	var left, right ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Shadows))),
			),
			Hand: ct.Cards(OratorHissaro),
		},
	})
	left.Exhaust()
	right.Exhaust()

	h.P1.Play(OratorHissaro)
	h.P1.ClickOption("Between") // deploy between the two neighbors

	if left.Exhausted() || right.Exhausted() {
		t.Error("neighbors should have been readied")
	}
	if left.AmberOn() != 1 || right.AmberOn() != 1 {
		t.Errorf("neighbors exalted = %d/%d, want 1/1", left.AmberOn(), right.AmberOn())
	}
	if got := h.Game().House(left.ID()); got != card.House.Saurian {
		t.Errorf("left house = %v, want Saurian", got)
	}
	if got := h.Game().House(right.ID()); got != card.House.Saurian {
		t.Errorf("right house = %v, want Saurian", got)
	}
}
