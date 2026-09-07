package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wild Spirit
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Reap: this creature captures 1 Æmber from your opponent."
func TestWildSpirit(t *testing.T) {
	t.Run("its host captures 1 Æmber when it reaps", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
						WildSpirit,
					),
				),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Reap(host)

		h.Expect(host).AmberOn(1)
		h.P2.ExpectAmber(1)
	})
}
