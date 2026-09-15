package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Valdr
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	Valdr deals +2 Damage while attacking an enemy creature on the flank.
func TestValdr(t *testing.T) {
	t.Run("deals +2 damage while attacking a flank creature", func(t *testing.T) {
		var flank ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Valdr)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&flank, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(10))),
			)},
		})

		h.P1.Fight(Valdr, flank)

		h.Expect(flank).Damage(8) // 6 power + 2 flank bonus
	})

	t.Run("deals no bonus against a creature that is not on a flank", func(t *testing.T) {
		var middle ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Valdr)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(10)),
				ct.Bind(&middle, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(10))),
				ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(10)),
			)},
		})

		h.P1.Fight(Valdr, middle)

		h.Expect(middle).Damage(6) // 6 power, no bonus
	})
}
