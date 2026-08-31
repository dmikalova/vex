package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Jehu the Bureaucrat
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	After you choose Sanctum as your active house, gain 2 Æmber.
func TestJehuTheBureaucrat(t *testing.T) {
	t.Run("gains 2 Æmber when Sanctum is chosen as the active house", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(JehuTheBureaucrat)},
		})

		// The opening house choice happens before cards are placed, so cycle round
		// to P1's next turn and choose Sanctum with Jehu in play.
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Sanctum)

		h.P1.ExpectAmber(2)
	})

	t.Run("does nothing when a different house is chosen", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(JehuTheBureaucrat)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)

		h.P1.ExpectAmber(0)
	})
}
