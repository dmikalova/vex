package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Too Much to Protect
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Steal all but 6 Æmber from your opponent.
func TestTooMuchToProtect(t *testing.T) {
	t.Run("steals the Æmber above six", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(TooMuchToProtect)},
			P2: ct.Side{Amber: 10},
		})

		h.P1.Play(TooMuchToProtect)

		h.P1.ExpectAmber(5) // 4 stolen, plus the Æmber bonus
		h.P2.ExpectAmber(6)
	})

	t.Run("steals nothing when the opponent has six or less", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(TooMuchToProtect)},
			P2: ct.Side{Amber: 6},
		})

		h.P1.Play(TooMuchToProtect)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(6)
	})
}
