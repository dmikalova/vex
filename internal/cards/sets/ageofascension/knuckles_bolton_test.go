package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Knuckles Bolton
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive, Skirmish.
func TestKnucklesBolton(t *testing.T) {
	t.Run("takes no damage back when fighting", func(t *testing.T) {
		var knuckles, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&knuckles, KnucklesBolton)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Fight(knuckles, foe)

		h.Expect(knuckles).Damage(0)
		h.Expect(foe).Damage(3)
	})
}
