package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Blood of Titans
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains +5 power.
func TestBloodOfTitans(t *testing.T) {
	t.Run("grants its host +5 power while attached", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
						BloodOfTitans,
					),
				),
			},
		})

		h.Expect(host).Power(9) // 4 base + 5
	})
}
