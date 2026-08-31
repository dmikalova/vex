package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Way of the Wolf
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains skirmish.
func TestWayOfTheWolf(t *testing.T) {
	t.Run("grants the host skirmish, sparing it return damage", func(t *testing.T) {
		var host, wall ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
						WayOfTheWolf,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&wall, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(10))),
				),
			},
		})

		h.P1.Fight(host, wall)

		h.Expect(host).
			Damage(0).
			At(ct.PlayArea)
		// skirmish: no return damage from the 10-power wall
	})
}
