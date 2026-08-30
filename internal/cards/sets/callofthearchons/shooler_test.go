package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shooler
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Play: If your opponent has 4 Æmber or more, steal 1 Æmber.
func TestShooler(t *testing.T) {
	t.Run("steals 1 Æmber when the opponent has 4 or more", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(Shooler)},
			P2: ct.Side{Amber: 4},
		})

		h.P1.Play(Shooler)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(3)
	})

	t.Run("does nothing when the opponent has fewer than 4", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(Shooler)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(Shooler)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
	})
}
