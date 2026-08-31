package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Way of the Bear
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains +2 assault.
func TestWayOfTheBear(t *testing.T) {
	t.Run("grants its host +2 assault before fight damage", func(t *testing.T) {
		var host, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
						WayOfTheBear,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(10))),
				),
			},
		})

		h.P1.Fight(host, foe)

		h.Expect(foe).Damage(5) // 2 assault + 3 fight
	})
}
