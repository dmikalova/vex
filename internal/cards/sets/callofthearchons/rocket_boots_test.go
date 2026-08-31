package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Rocket Boots
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "Fight/Reap: If this is the first time this creature was used this turn, ready it."
func TestRocketBoots(t *testing.T) {
	t.Run("readies the host only the first time it is used each turn", func(t *testing.T) {
		var host, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(5))),
						RocketBoots,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.Reap(host)
		h.Expect(host).Ready()

		h.P1.Fight(host, foe)
		h.Expect(host).Exhausted()

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)

		h.P1.Reap(host)
		h.Expect(host).Ready()
	})
}
