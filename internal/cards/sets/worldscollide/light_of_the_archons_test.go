package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Light of the Archons
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature gains +1 power and +1 armor for each upgrade attached to it.
func TestLightOfTheArchons(t *testing.T) {
	t.Run("with one upgrade attached the host gains +1 power and +1 armor", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.Power(4))),
						LightOfTheArchons,
					),
				),
			},
		})

		h.Expect(host).Power(5).Armor(1)
	})

	t.Run("the bonus scales with every upgrade on the host", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.Power(4))),
						LightOfTheArchons,
						ct.Upgrade(ct.AemberBonus(1)),
					),
				),
			},
		})

		// Two upgrades on the host: +2 power and +2 armor.
		h.Expect(host).Power(6).Armor(2)
	})
}
