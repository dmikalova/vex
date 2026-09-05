package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cull the Weak
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Destroy the least powerful enemy creature.
func TestCullTheWeak(t *testing.T) {
	t.Run("destroys the least powerful enemy creature", func(t *testing.T) {
		var weak, strong ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(CullTheWeak)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&weak, ct.Creature(ct.Power(2))),
				ct.Bind(&strong, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(CullTheWeak)

		h.Expect(weak).At(ct.Discard)
		h.Expect(strong).At(ct.PlayArea)
	})
}
