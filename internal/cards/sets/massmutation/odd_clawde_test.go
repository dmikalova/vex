package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Odd Clawde
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Mutant • Scientist
//
//	Action: If your opponent has an odd amount of Æmber, steal 1 Æmber.
func TestOddClawde(t *testing.T) {
	t.Run("steals when the opponent has an odd amount of Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(OddClawde)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(OddClawde)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("does nothing when the opponent has an even amount of Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(OddClawde)},
			P2: ct.Side{Amber: 2},
		})

		h.P1.UseAction(OddClawde)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(2)
	})
}
