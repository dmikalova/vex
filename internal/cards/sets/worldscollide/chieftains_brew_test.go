package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Chieftain's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight: Ready and fight with a neighboring creature."
func TestChieftainsBrew(t *testing.T) {
	t.Run("host fighting readies and fights a neighboring creature", func(t *testing.T) {
		var host, neighbor, weakFoe, bigFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
						ChieftainsBrew,
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&weakFoe, ct.Creature(ct.Power(1))),
				ct.Bind(&bigFoe, ct.Creature(ct.Power(10))),
			)},
		})

		h.P1.Fight(host, weakFoe)

		h.Expect(bigFoe).Damage(4)
	})
}
