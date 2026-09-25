package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Pestering Blow
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to a creature and enrage it.
func TestPesteringBlow(t *testing.T) {
	t.Run("deals 1 damage to a creature and enrages it", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(PesteringBlow),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4))))},
		})

		h.P1.Play(PesteringBlow)

		h.Expect(foe).Damage(1)
		if !h.Game().Enraged(foe.ID()) {
			t.Error("the creature should be enraged")
		}
	})
}
