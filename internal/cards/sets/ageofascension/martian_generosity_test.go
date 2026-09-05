package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Martian Generosity
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Lose all your Æmber, and for each Æmber you lost this way, draw 2 cards.
func TestMartianGenerosity(t *testing.T) {
	t.Run("loses all aember and draws 2 cards for each aember lost", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Amber: 2,
				Hand:  ct.Cards(MartianGenerosity),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.Power(1))),
					ct.Creature(ct.Power(1)),
					ct.Creature(ct.Power(1)),
					ct.Creature(ct.Power(1)),
					ct.Creature(ct.Power(1)),
					ct.Creature(ct.Power(1)),
				),
			},
		})

		h.P1.Play(MartianGenerosity)

		h.P1.ExpectAmber(0)
		h.Expect(top).At(ct.Hand)
	})
}
