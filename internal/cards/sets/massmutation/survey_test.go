package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Survey
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Look at the top 2 cards of your deck and discard 1.
//	Enhance Draw.
func TestSurvey(t *testing.T) {
	t.Run("looks at the top 2 cards and discards 1, the other stays on top", func(t *testing.T) {
		var top, second, third ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(Survey),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.Power(3))),
					ct.Bind(&second, ct.Creature(ct.Power(4))),
					ct.Bind(&third, ct.Creature(ct.Power(5))),
				),
			},
		})

		h.P1.Play(Survey)
		h.P1.ClickCard(top) // discard one of the two shown

		h.Expect(top).At(ct.Discard)
		h.Expect(second).At(ct.Deck)
		h.Expect(third).At(ct.Deck)
	})
}
