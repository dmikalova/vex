package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ghostly Hand
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  2
//
//	Play: If your opponent has exactly 1 Æmber, steal 1 Æmber.
func TestGhostlyHand(t *testing.T) {
	t.Run("gains its bonus and steals 1 at exactly 1 opponent Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(GhostlyHand)},
			P2: ct.Side{Amber: 1},
		})

		h.P1.Play(GhostlyHand)

		h.P1.ExpectAmber(3) // 2 bonus + 1 stolen
		h.P2.ExpectAmber(0)
	})

	t.Run("only gains its bonus when the opponent is not at exactly 1", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(GhostlyHand)},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Play(GhostlyHand)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(2)
	})
}
