package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Finishing Blow
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a damaged creature -> steal 1 Æmber.
func TestFinishingBlow(t *testing.T) {
	t.Run("destroys a damaged creature and steals 1 Æmber", func(t *testing.T) {
		var dmg ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(FinishingBlow)},
			P2: ct.Side{
				Amber: 3,
				InPlay: ct.Cards(
					ct.Bind(&dmg, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})
		dmg.Damaged(2)

		h.P1.Play(FinishingBlow)

		h.Expect(dmg).At(ct.Discard)
		h.P1.ExpectAmber(2) // 1 Æmber bonus pip + 1 stolen
		h.P2.ExpectAmber(2) // 3 - 1 stolen
	})
}
