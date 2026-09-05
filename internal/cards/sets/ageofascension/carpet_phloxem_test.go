package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Carpet Phloxem
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If there are no friendly creatures in play, deal 4 damage to each creature.
func TestCarpetPhloxem(t *testing.T) {
	t.Run("deals 4 damage to each creature when you control no creatures", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(CarpetPhloxem)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Play(CarpetPhloxem)

		h.Expect(foe).Damage(4)
	})
}
