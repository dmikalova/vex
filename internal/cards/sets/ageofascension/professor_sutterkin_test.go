package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Professor Sutterkin
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Scientist
//
//	Reap: For each friendly Logos creature, draw a card.
func TestProfessorSutterkin(t *testing.T) {
	t.Run("draws a card for each friendly Logos creature", func(t *testing.T) {
		var top1, top2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ProfessorSutterkin,
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
				),
				Deck: ct.Cards(
					ct.Bind(&top1, ct.Creature(ct.Power(3))),
					ct.Bind(&top2, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.Reap(ProfessorSutterkin)

		h.Expect(top1).At(ct.Hand)
		h.Expect(top2).At(ct.Hand)
	})
}
