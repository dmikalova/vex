package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Drumble
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Play: If your opponent has 7 Æmber or more, Drumble captures all your opponent's Æmber.
func TestDrumble(t *testing.T) {
	t.Run("captures the opponent's whole pool at 7+ Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Drumble),
			},
			P2: ct.Side{Amber: 7},
		})

		h.P1.Play(Drumble)

		h.Expect(Drumble).AmberOn(7)
		h.P2.ExpectAmber(0)
	})

	t.Run("captures nothing below the threshold", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Drumble),
			},
			P2: ct.Side{Amber: 6},
		})

		h.P1.Play(Drumble)

		h.Expect(Drumble).AmberOn(0)
		h.P2.ExpectAmber(6) // untouched
	})
}
