package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mega Ganger Chieftain
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Giant
//
//	Play: Ready and fight with a neighboring creature.
func TestMegaGangerChieftain(t *testing.T) {
	t.Run("readies and fights with a neighboring creature", func(t *testing.T) {
		var neighbor, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(MegaGangerChieftain),
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(1))),
				),
			},
		})

		h.P1.Play(MegaGangerChieftain)

		h.Expect(foe).At(ct.Discard)
		h.Expect(neighbor).Exhausted()
	})
}
