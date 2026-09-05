package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cowfyne
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Before Fight: Deal 2 damage to each neighbor of the creature Cowfyne fights.
func TestCowfyne(t *testing.T) {
	t.Run("deals 2 damage to each neighbor of the creature it fights", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Cowfyne),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(20))),
					ct.Bind(&middle, ct.Creature(ct.Power(20))),
					ct.Bind(&right, ct.Creature(ct.Power(20))),
				),
			},
		})

		h.P1.Fight(Cowfyne, middle)

		h.Expect(left).Damage(2)
		h.Expect(right).Damage(2)
	})
}
