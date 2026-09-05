package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// "Lion" Bautrem
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Human • Knight
//
//	Deploy.
//	Each neighboring creature gains +2 power.
func TestLionBautrem(t *testing.T) {
	t.Run("gives each battleline neighbor +2 power but not distant creatures", func(t *testing.T) {
		var left, right, far ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
					LionBautrem,
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
					ct.Bind(&far, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
		})

		h.Expect(left).Power(5)
		h.Expect(right).Power(5)
		h.Expect(far).Power(3) // not a neighbor
	})
}
