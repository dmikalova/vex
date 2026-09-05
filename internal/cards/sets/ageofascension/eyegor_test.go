package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Eyegor
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg
//
//	Play: Look at the top 3 cards of your deck, put 1 into your hand, and discard the others.
func TestEyegor(t *testing.T) {
	t.Run("keeps the chosen card and discards the others", func(t *testing.T) {
		var top, middle, third, bottom ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(Eyegor),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.Power(3))),
					ct.Bind(&middle, ct.Creature(ct.Power(3))),
					ct.Bind(&third, ct.Creature(ct.Power(3))),
					ct.Bind(&bottom, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.Play(Eyegor)
		h.P1.ClickCard(middle) // keep one of the top three

		h.Expect(middle).At(ct.Hand)
		h.Expect(top).At(ct.Discard)
		h.Expect(third).At(ct.Discard)
		h.Expect(bottom).At(ct.Deck)
	})
}
