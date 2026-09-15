package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// City-State Interest
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Each friendly creature captures 1 Æmber from your opponent.
func TestCityStateInterest(t *testing.T) {
	t.Run("each friendly creature captures 1 Æmber from your opponent", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(CityStateInterest),
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Saurian))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Saurian))),
				),
			},
			P2: ct.Side{
				Amber: 3,
			},
		})

		h.P1.Play(CityStateInterest)

		h.Expect(a).AmberOn(1)
		h.Expect(b).AmberOn(1)
		h.P2.ExpectAmber(1)
	})
}
