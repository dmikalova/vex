package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Binate Rupture
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Alpha.
//	Play: For each Æmber in your pool, gain 1 Æmber, and for each Æmber in your opponent's pool, your opponent gains 1 Æmber.
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
