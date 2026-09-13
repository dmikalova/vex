package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mega Cowfyne
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  7
//	Traits: Giant
//
//	Before Fight: Deal 2 damage to each neighbor of the Creature Mega Cowfyne fights.
func TestMegaCowfyne(t *testing.T) {
	t.Run("deals 2 damage to each neighbor of the creature it fights", func(t *testing.T) {
		var target, left, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(MegaCowfyne)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(4))),
				ct.Bind(&target, ct.Creature(ct.Power(4))),
				ct.Bind(&right, ct.Creature(ct.Power(4))),
			)},
		})

		h.P1.Fight(MegaCowfyne, target)

		h.Expect(left).Damage(2)
		h.Expect(right).Damage(2)
	})
}
