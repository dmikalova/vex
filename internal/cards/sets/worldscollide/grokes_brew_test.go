package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Groke's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight: Your opponent loses 1 Æmber."
func TestGrokesBrew(t *testing.T) {
	t.Run("opponent loses 1 Æmber when the host fights", func(t *testing.T) {
		var host, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
						GrokesBrew,
					),
				),
			},
			P2: ct.Side{
				Amber:  3,
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Fight(host, foe)

		h.P2.ExpectAmber(2)
	})
}
