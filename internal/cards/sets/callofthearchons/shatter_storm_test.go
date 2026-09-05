package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shatter Storm
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Lose all your Æmber, and for each Æmber you lost this way, your opponent loses 3 Æmber.
func TestShatterStorm(t *testing.T) {
	t.Run("empties your pool and drains triple that from your opponent", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ShatterStorm),
				Amber: 3,
			},
			P2: ct.Side{Amber: 10},
		})

		h.P1.Play(ShatterStorm)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(1) // 10 - 3*3
	})

	t.Run("drains nothing when you have no Æmber to lose", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ShatterStorm),
			},
			P2: ct.Side{Amber: 4},
		})

		h.P1.Play(ShatterStorm)

		h.P2.ExpectAmber(4)
	})
}
