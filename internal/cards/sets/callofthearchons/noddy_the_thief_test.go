package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Noddy the Thief
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Steal 1 Æmber.
func TestNoddyTheThief(t *testing.T) {
	t.Run("steals 1 Æmber as an action", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(NoddyTheThief)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(NoddyTheThief)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})
}
