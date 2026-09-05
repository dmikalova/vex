package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Whistling Darts
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 1 damage to each enemy creature.
func TestWhistlingDarts(t *testing.T) {
	t.Run("deals 1 damage to each enemy creature when played", func(t *testing.T) {
		var foe1, foe2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(WhistlingDarts)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe1, ct.Creature(ct.Power(20))),
					ct.Bind(&foe2, ct.Creature(ct.Power(20))),
				),
			},
		})

		h.P1.Play(WhistlingDarts)

		h.Expect(foe1).Damage(1)
		h.Expect(foe2).Damage(1)
	})
}
