package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ganger Chieftain
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Play: Ready and fight with a neighboring creature.
func TestGangerChieftain(t *testing.T) {
	t.Run("readies and fights with a neighboring creature when played", func(t *testing.T) {
		var neighbor, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
				Hand: ct.Cards(GangerChieftain),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Play(GangerChieftain)

		h.Expect(foe).At(ct.Discard)   // the neighbor fought and destroyed it
		h.Expect(neighbor).Exhausted() // exhausted after fighting
	})
}
