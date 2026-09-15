package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Containment Field
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "After this creature is used, destroy this creature."
func TestContainmentField(t *testing.T) {
	t.Run("destroys its host after the host reaps", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
						ContainmentField,
					),
				),
			},
		})

		h.P1.Reap(host)

		h.Expect(host).At(ct.Discard)
	})

	t.Run("destroys its host after the host fights", func(t *testing.T) {
		var host, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
						ContainmentField,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.Fight(host, enemy)

		h.Expect(host).At(ct.Discard)
	})
}
