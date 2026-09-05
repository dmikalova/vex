package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Flamethrower
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Action: Deal 1 damage to a creature that is not on a flank and 1 damage to each of its neighbors.
func TestFlamethrower(t *testing.T) {
	t.Run("action deals 1 damage to a creature and 1 to each neighbor", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Flamethrower)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(20))),
					ct.Bind(&middle, ct.Creature(ct.Power(20))),
					ct.Bind(&right, ct.Creature(ct.Power(20))),
				),
			},
		})

		h.P1.UseAction(Flamethrower)

		h.Expect(middle).Damage(1)
		h.Expect(left).Damage(1)
		h.Expect(right).Damage(1)
	})
}
