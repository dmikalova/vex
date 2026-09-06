package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nightforge
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If you have not forged a key this turn, you may forge a key at +4 Æmber current cost.
func TestNightforge(t *testing.T) {
	t.Run("may forge a key at +4 current cost", func(t *testing.T) {
		var nf ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 20,
				Hand:  ct.Cards(ct.Bind(&nf, Nightforge)),
			},
		})

		h.P1.Play(nf)
		h.P1.ClickOption("Yes")

		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(11)
	})

	t.Run("may decline the forge", func(t *testing.T) {
		var nf ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 20,
				Hand:  ct.Cards(ct.Bind(&nf, Nightforge)),
			},
		})

		h.P1.Play(nf)
		h.P1.ClickOption("No")

		h.P1.ExpectKeys(0)
		h.P1.ExpectAmber(21)
	})
}
