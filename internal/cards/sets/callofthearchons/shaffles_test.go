package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shaffles
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	At the end of your turn, your opponent loses 1 Æmber.
func TestShaffles(t *testing.T) {
	t.Run("drains 1 Æmber from the opponent at the end of the turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Shaffles)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.EndTurn()

		h.P2.ExpectAmber(2)
	})
}
