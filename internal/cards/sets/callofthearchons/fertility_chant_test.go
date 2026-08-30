package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fertility Chant
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  4
//
//	Play: Your opponent gains 2 Æmber.
func TestFertilityChant(t *testing.T) {
	t.Run("gives the controller 4 Æmber pips and the opponent 2", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(FertilityChant)},
		})

		h.P1.Play(FertilityChant)

		h.P1.ExpectAmber(4)
		h.P2.ExpectAmber(2)
	})
}
