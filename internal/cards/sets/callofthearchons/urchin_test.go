package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Urchin
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Play: Steal 1 Æmber.
func TestUrchin(t *testing.T) {
	t.Run("steals 1 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Urchin),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Play(Urchin)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
