package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Treasure Map
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If you played exactly 1 card this turn, gain 3 Æmber, and you cannot play cards for the remainder of the turn.
func TestTreasureMap(t *testing.T) {
	t.Run(
		"pays out when it is the first card played, then shuts the turn down",
		func(t *testing.T) {
			var treasureMap, other ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Shadows,
					Hand: ct.Cards(
						ct.Bind(&treasureMap, TreasureMap),
						ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
					),
				},
			})

			h.P1.Play(treasureMap)

			// 1 Æmber bonus plus the 3 it pays out.
			h.P1.ExpectAmber(4)
			h.P1.ExpectCannotPlay(other)
		},
	)

	t.Run("pays out nothing when another card came first", func(t *testing.T) {
		var treasureMap, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand: ct.Cards(
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
					ct.Bind(&treasureMap, TreasureMap),
				),
			},
		})

		h.P1.Play(other)
		h.P1.Play(treasureMap)

		h.P1.ExpectAmber(1)
	})
}
