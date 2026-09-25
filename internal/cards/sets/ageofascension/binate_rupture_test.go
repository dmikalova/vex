package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Binate Rupture
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Alpha.
//	Play: Each player gains Æmber equal to the Æmber in their pool.
func TestBinateRupture(t *testing.T) {
	t.Run("each player gains aember equal to the aember in their pool", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Amber: 3,
				Hand:  ct.Cards(BinateRupture),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Play(BinateRupture)

		h.P1.ExpectAmber(6)
		h.P2.ExpectAmber(4)
	})
}
