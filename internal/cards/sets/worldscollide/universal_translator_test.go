package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Universal Translator
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Use a non-Star Alliance creature."
func TestUniversalTranslator(t *testing.T) {
	t.Run(
		"its host may use a friendly non-Star Alliance creature when it reaps",
		func(t *testing.T) {
			var host, ally ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.StarAlliance,
					InPlay: ct.Cards(
						ct.Upgraded(
							ct.Bind(
								&host,
								ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
							),
							UniversalTranslator,
						),
						ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					),
				},
			})

			h.P1.Reap(host)

			h.Expect(ally).Exhausted()
			h.P1.ExpectAmber(2)
		},
	)
}
