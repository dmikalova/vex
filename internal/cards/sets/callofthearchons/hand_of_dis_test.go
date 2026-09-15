package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hand of Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy a creature that is not on a flank.
func TestHandOfDis(t *testing.T) {
	t.Run("destroys a chosen creature that is not on a flank", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(HandOfDis)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&mid, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Play(HandOfDis) // mid is the only non-flank creature, so it is the forced target

		h.Expect(mid).At(ct.Discard)
		h.Expect(left).At(ct.PlayArea)
		h.Expect(right).At(ct.PlayArea)
	})
}
