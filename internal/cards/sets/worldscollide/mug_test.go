package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mug
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a creature - move 1 Æmber from it to your pool. Deal 2 damage to it.
func TestMug(t *testing.T) {
	t.Run("takes 1 Æmber from a creature and deals 2 damage to it", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(Mug)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5))))},
		})
		h.Game().State.Cards[foe.ID()].Amber = 2

		h.P1.Play(Mug)

		h.Expect(foe).AmberOn(1).Damage(2)
		h.P1.ExpectAmber(2) // 1 from Mug's bonus, 1 taken from the creature
	})
}
