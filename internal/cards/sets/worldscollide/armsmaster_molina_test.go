package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Armsmaster Molina
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Hazardous 3.
//	Each neighboring creature gains hazardous 3.
func TestArmsmasterMolina(t *testing.T) {
	var neighbor ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.StarAlliance,
			InPlay: ct.Cards(
				ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3))),
				ArmsmasterMolina,
			),
		},
	})

	if got := h.Game().Hazardous(neighbor.ID()); got != 3 {
		t.Errorf("neighbor hazardous = %d, want 3", got)
	}
}
