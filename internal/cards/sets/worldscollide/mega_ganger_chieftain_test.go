package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mega Ganger Chieftain
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  7
//	Traits: Giant
//
//	Play: You may ready and fight with a neighboring Creature.
func TestMegaGangerChieftain(t *testing.T) {
	t.Run("may ready and fight with a neighboring creature", func(t *testing.T) {
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
		h.P1.ClickCard(neighbor)

		h.Expect(foe).At(ct.Discard)
		h.Expect(neighbor).Exhausted()
	})
}
