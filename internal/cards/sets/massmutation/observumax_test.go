package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Observe-u-Max
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: This creature captures 1 Æmber from your opponent."
func TestObservuMax(t *testing.T) {
	t.Run("host captures 1 Æmber when it reaps", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
						ObservuMax,
					),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(host)

		h.Expect(host).AmberOn(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("host captures 1 Æmber when it fights", func(t *testing.T) {
		var host, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(
							ct.OfHouse(card.House.StarAlliance),
							ct.Power(4),
						)),
						ObservuMax,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
				Amber: 3,
			},
		})

		h.P1.Fight(host, enemy)

		h.Expect(host).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
