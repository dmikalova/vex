package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// They're Everywhere!
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to each enemy flank creature. Deal 1 damage to each enemy creature that is not on a flank.
func TestTheyreEverywhere(t *testing.T) {
	t.Run("deals 2 damage to flank enemies and 1 damage to the rest", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(TheyreEverywhere)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(20))),
					ct.Bind(&middle, ct.Creature(ct.Power(20))),
					ct.Bind(&right, ct.Creature(ct.Power(20))),
				),
			},
		})

		h.P1.Play(TheyreEverywhere)

		h.Expect(left).Damage(2)
		h.Expect(right).Damage(2)
		h.Expect(middle).Damage(1)
	})
}
