package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dimension Door
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For the remainder of the turn, instead of gaining Æmber from reaping, steal the same amount.
func TestDimensionDoor(t *testing.T) {
	t.Run("reaps steal Æmber from the opponent instead of gaining it", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(DimensionDoor),
				InPlay: ct.Cards(
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(DimensionDoor)
		h.P1.Reap(reaper)

		h.P1.ExpectAmber(1) // stole 1 rather than gaining from the supply
		h.P2.ExpectAmber(2)
	})
}
