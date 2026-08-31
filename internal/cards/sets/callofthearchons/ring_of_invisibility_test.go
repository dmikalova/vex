package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ring of Invisibility
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains elusive and skirmish.
func TestRingOfInvisibility(t *testing.T) {
	t.Run("grants skirmish so the host takes no retaliation when it fights", func(t *testing.T) {
		var host, wall ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(4))),
						RingOfInvisibility,
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

		h.Expect(host).Damage(0) // granted skirmish
	})
}
