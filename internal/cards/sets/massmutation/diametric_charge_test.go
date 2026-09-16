package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Diametric Charge
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to a creature and 2 damage to each of its neighbors.
func TestDiametricCharge(t *testing.T) {
	t.Run("deals 1 to a creature and 2 to each neighbor", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(DiametricCharge)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(5))),
				ct.Bind(&mid, ct.Creature(ct.Power(5))),
				ct.Bind(&right, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(DiametricCharge)
		h.P1.ClickCard(mid)

		h.Expect(mid).Damage(1)
		h.Expect(left).Damage(2)
		h.Expect(right).Damage(2)
	})
}
