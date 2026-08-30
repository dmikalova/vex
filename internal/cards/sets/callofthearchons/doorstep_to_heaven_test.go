package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Doorstep to Heaven
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Each player with 6 Æmber or more loses all but 5 Æmber.
func TestDoorstepToHeaven(t *testing.T) {
	t.Run("reduces each player with 6+ Æmber to 5", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(DoorstepToHeaven), Amber: 4},
			P2: ct.Side{Amber: 9},
		})

		h.P1.Play(DoorstepToHeaven)

		// Playing the card grants its 1 Æmber pip (4 -> 5), which is not over the
		// cap, so P1 stays at 5; P2 is reduced from 9 to 5.
		h.P1.ExpectAmber(5)
		h.P2.ExpectAmber(5)
	})
}
