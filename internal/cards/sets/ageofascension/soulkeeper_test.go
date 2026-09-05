package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Soulkeeper
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Destroyed: Destroy the 1 most powerful enemy creatures."
func TestSoulkeeper(t *testing.T) {
	t.Run(
		"destroys the most powerful enemy creature when its host is destroyed",
		func(t *testing.T) {
			var host, weak, strong ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					InPlay: ct.Cards(
						ct.Upgraded(
							ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(2))),
							Soulkeeper,
						),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&weak, ct.Creature(ct.Power(6))),
						ct.Bind(&strong, ct.Creature(ct.Power(10))),
					),
				},
			})

			h.P1.Fight(host, weak)

			h.Expect(host).At(ct.Discard)
			h.Expect(strong).At(ct.Discard)
		},
	)
}
