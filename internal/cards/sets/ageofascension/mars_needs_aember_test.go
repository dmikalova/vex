package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mars Needs Aember
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Each enemy damaged non-Mars creature captures 1 Æmber from your opponent.
func TestMarsNeedsAember(t *testing.T) {
	t.Run("each damaged enemy non-Mars creature captures 1 aember", func(t *testing.T) {
		var damaged ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(MarsNeedsAember),
			},
			P2: ct.Side{
				Amber: 3,
				InPlay: ct.Cards(
					ct.Bind(&damaged, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(5))),
				),
			},
		})
		damaged.Damaged(1)

		h.P1.Play(MarsNeedsAember)

		h.Expect(damaged).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
