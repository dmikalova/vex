package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Duskrunner
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "Reap: Steal 1 Æmber."
func TestDuskrunner(t *testing.T) {
	t.Run("grants the host Reap: Steal 1 Æmber", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
						Duskrunner,
					),
				),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Reap(host)

		h.P1.ExpectAmber(2) // 1 reap + 1 stolen
		h.P2.ExpectAmber(1)
	})
}
