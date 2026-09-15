package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Warriors' Refrain
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Stun each creature with power 3 or lower.
func TestWarriorsRefrain(t *testing.T) {
	t.Run("stuns each creature with power 3 or lower", func(t *testing.T) {
		var weak, strong ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(WarriorsRefrain)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&weak, ct.Creature(ct.Power(3))),
				ct.Bind(&strong, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(WarriorsRefrain)

		h.Expect(weak).Stunned(true)
		h.Expect(strong).Stunned(false)
	})
}
