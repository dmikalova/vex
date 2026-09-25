package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Golden Aura
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a creature. Fully heal it. For the remainder of the turn, it belongs to house Sanctum and cannot be dealt damage.
func TestGoldenAura(t *testing.T) {
	t.Run("fully heals the chosen creature", func(t *testing.T) {
		var target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(GoldenAura),
				InPlay: ct.Cards(
					ct.Bind(&target, ct.Creature(ct.Power(5))),
				),
			},
		})
		target.Damaged(3)

		h.P1.Play(GoldenAura)

		h.Expect(target).Damage(0)
	})
}
