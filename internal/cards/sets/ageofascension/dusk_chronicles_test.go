package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dusk Chronicles
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If your opponent has more Æmber than you, draw a card. If you have more Æmber than your opponent, archive a card from your hand.
func TestDuskChronicles(t *testing.T) {
	t.Run("draws when opponent has more aember", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 0, // 1 after the play bonus
				Hand:  ct.Cards(DuskChronicles),
				Deck:  ct.Cards(ct.Bind(&top, ct.Creature(ct.Power(1)))),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(DuskChronicles)

		h.Expect(top).At(ct.Hand)
	})

	t.Run("archives when you have more aember", func(t *testing.T) {
		var spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 5, // 6 after the play bonus
				Hand: ct.Cards(
					DuskChronicles,
					ct.Bind(&spare, ct.Creature(ct.Power(1))),
				),
			},
			P2: ct.Side{Amber: 1},
		})

		h.P1.Play(DuskChronicles)

		h.Expect(spare).At(ct.Archives)
	})
}
