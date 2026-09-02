package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Total Recall
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: For each friendly ready creature in play, gain 1 Æmber. Put each friendly creature into its owner's hand.
func TestTotalRecall(t *testing.T) {
	t.Run("gains 1 Æmber per ready friendly creature and returns them all", func(t *testing.T) {
		var spent, ready ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(TotalRecall),
				InPlay: ct.Cards(
					ct.Bind(&spent, ct.Creature(ct.OfHouse(card.House.Mars))),
					ct.Bind(&ready, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})

		h.P1.Reap(spent)
		h.P1.Play(TotalRecall)

		h.P1.ExpectAmber(3) // 1 reaped, 1 Æmber bonus, 1 for the ready creature
		h.Expect(spent).At(ct.Hand)
		h.Expect(ready).At(ct.Hand)
	})
}
