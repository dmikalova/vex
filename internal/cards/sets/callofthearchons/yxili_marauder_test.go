package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Yxili Marauder
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Martian • Soldier
//
//	Yxili Marauder gains +1 power for each Æmber on it.
//	Play: For each friendly ready Mars creature, Yxili Marauder captures 1 Æmber from your opponent.
func TestYxiliMarauder(t *testing.T) {
	t.Run(
		"captures 1 Æmber per friendly ready Mars creature and grows with it",
		func(t *testing.T) {
			var yxili, ready, exhausted ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Mars,
					Hand:  ct.Cards(ct.Bind(&yxili, YxiliMarauder)),
					InPlay: ct.Cards(
						ct.Bind(&ready, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
						ct.Bind(&exhausted, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					),
				},
				P2: ct.Side{Amber: 5},
			})
			exhausted.Exhaust()

			h.P1.Play(YxiliMarauder)

			// Yxili enters play exhausted, so only the other ready Mars creature counts.
			h.Expect(yxili).AmberOn(1)
			h.P2.ExpectAmber(4)
			h.Expect(yxili).Power(3)
		},
	)
}
