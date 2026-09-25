package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mars First
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Ready and use a friendly Mars creature.
func TestMarsFirst(t *testing.T) {
	t.Run("readies and uses a friendly Mars creature", func(t *testing.T) {
		var trooper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(MarsFirst),
				InPlay: ct.Cards(
					ct.Bind(&trooper, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				),
			},
			P2: ct.Side{},
		})
		trooper.Exhaust()

		h.P1.Play(MarsFirst)

		h.P1.ExpectAmber(2)
	})
}
