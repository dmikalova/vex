package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pit Demon
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Action: Steal 1 Æmber.
func TestPitDemon(t *testing.T) {
	t.Run("steals 1 Æmber as an action", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(PitDemon)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(PitDemon)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})
}
