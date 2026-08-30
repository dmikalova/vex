package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Radiant Truth
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Stun each enemy creature that is not on a flank.
func TestRadiantTruth(t *testing.T) {
	t.Run("stuns each enemy creature that is not on a flank", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(RadiantTruth)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&mid, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Play(RadiantTruth)

		h.Expect(mid).Stunned(true)   // the interior creature is stunned
		h.Expect(left).Stunned(false) // flanks are spared
		h.Expect(right).Stunned(false)
	})
}
