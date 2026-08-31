package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Krump
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	After a creature is destroyed fighting Krump, your opponent loses 1 Æmber.
func TestKrump(t *testing.T) {
	t.Run("opponent loses 1 Æmber when a creature is destroyed fighting Krump", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Krump)},
			P2: ct.Side{
				Amber: 3,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Fight(Krump, foe)

		h.Expect(foe).At(ct.Discard)
		h.P2.ExpectAmber(2)
	})
}
