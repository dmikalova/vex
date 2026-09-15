package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Slimy Jark
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Goblin
//
//	Skirmish, Elusive.
//	Fight: Enrage the creature Slimy Jark fought.
func TestSlimyJark(t *testing.T) {
	t.Run("enrages the creature it fights and takes no damage back", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(SlimyJark)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(SlimyJark, foe)

		h.Expect(foe).Damage(1)
		if !h.Game().Enraged(foe.ID()) {
			t.Errorf("%s should be enraged", foe.Name())
		}
		h.Expect(SlimyJark).At(ct.PlayArea).Damage(0)
	})
}
