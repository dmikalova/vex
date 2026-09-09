package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Brobnar Plant
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	After a player chooses Brobnar as their active house, gain 1 Æmber.
func TestBrobnarPlant(t *testing.T) {
	t.Run("gains 1 Æmber when a player chooses Brobnar", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(BrobnarPlant)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P1.ExpectAmber(1)
	})

	t.Run("does nothing when a different house is chosen", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(BrobnarPlant)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Logos)

		h.P1.ExpectAmber(0)
	})
}
