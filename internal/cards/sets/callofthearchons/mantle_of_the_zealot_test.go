package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mantle of the Zealot
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains versatile.
func TestMantleOfTheZealot(t *testing.T) {
	t.Run("grants versatile so an out-of-house host can reap", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
						MantleOfTheZealot,
					),
				),
			},
		})

		h.P1.Reap(host) // versatile lets the Dis creature reap on a Sanctum turn

		h.P1.ExpectAmber(1)
		h.Expect(host).Exhausted()
	})
}
