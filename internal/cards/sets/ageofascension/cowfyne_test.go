package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Cowfyne
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Splash-attack 2.
func TestCowfyne(t *testing.T) {
	t.Run("deals 2 splash damage to each neighbor while it fights", func(t *testing.T) {
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
		h.Expect(middle).Damage(5)
		h.Expect(right).Damage(2)
	})
}
