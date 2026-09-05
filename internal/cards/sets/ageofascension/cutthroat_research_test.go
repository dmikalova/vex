package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cutthroat Research
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If your opponent has 8 Æmber or more, steal 2 Æmber.
func TestCutthroatResearch(t *testing.T) {
	t.Run("steals 2 when the opponent has 8 Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(CutthroatResearch)},
			P2: ct.Side{Amber: 8},
		})

		h.P1.Play(CutthroatResearch)

		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(6)
	})

	t.Run("steals nothing below 8 Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(CutthroatResearch)},
			P2: ct.Side{Amber: 7},
		})

		h.P1.Play(CutthroatResearch)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(7)
	})
}
