package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Furtive Investors
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If your opponent has more Æmber than you, for each forged key your opponent has, gain 1 Æmber.
func TestFurtiveInvestors(t *testing.T) {
	t.Run("gains 1 aember per opponent key when the opponent has more aember", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(FurtiveInvestors)},
			P2: ct.Side{Amber: 5, Keys: 2},
		})

		h.P1.Play(FurtiveInvestors)

		h.P1.ExpectAmber(3)
	})
}
