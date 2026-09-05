package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Phloxem Spike
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If there are no friendly creatures in play, destroy each creature that is not on a flank.
func TestPhloxemSpike(t *testing.T) {
	t.Run("destroys each non-flank creature when you control no creatures", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(PhloxemSpike)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(20))),
					ct.Bind(&middle, ct.Creature(ct.Power(20))),
					ct.Bind(&right, ct.Creature(ct.Power(20))),
				),
			},
		})

		h.P1.Play(PhloxemSpike)

		h.Expect(middle).At(ct.Discard)
		h.Expect(left).At(ct.PlayArea)
		h.Expect(right).At(ct.PlayArea)
	})
}
