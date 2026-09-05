package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bingle Bangbang
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Goblin
//
//	Before Fight: Deal 5 damage to each neighbor of the creature Bingle Bangbang fights.
func TestBingleBangbang(t *testing.T) {
	t.Run("deals 5 damage to each neighbor of the creature it fights", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(BingleBangbang),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(20))),
					ct.Bind(&middle, ct.Creature(ct.Power(20))),
					ct.Bind(&right, ct.Creature(ct.Power(20))),
				),
			},
		})

		h.P1.Fight(BingleBangbang, middle)

		h.Expect(left).Damage(5)
		h.Expect(right).Damage(5)
	})
}
