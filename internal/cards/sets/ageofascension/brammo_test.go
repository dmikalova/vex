package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Brammo
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Giant • Knight
//
//	Play: Deal 2 damage to each enemy flank creature.
func TestBrammo(t *testing.T) {
	t.Run("deals 2 damage to each enemy flank creature", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(Brammo)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(5))),
				ct.Bind(&middle, ct.Creature(ct.Power(5))),
				ct.Bind(&right, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(Brammo)

		h.Expect(left).Damage(2)
		h.Expect(middle).Damage(0)
		h.Expect(right).Damage(2)
	})
}
