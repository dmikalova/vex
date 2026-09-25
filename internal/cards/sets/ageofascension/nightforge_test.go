package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Nightforge
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If you have not forged a key this turn, forge a key at +4 Æmber current cost -> purge Nightforge.
func TestNightforge(t *testing.T) {
	t.Run("forges a key at +4 current cost and purges itself", func(t *testing.T) {
		var nf ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 20,
				Hand:  ct.Cards(ct.Bind(&nf, Nightforge)),
			},
		})

		h.P1.Play(nf)

		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(11)
		h.Expect(nf).At(ct.Purge)
	})

	t.Run("stays in the discard pile when the forge is unaffordable", func(t *testing.T) {
		var nf ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Amber: 5,
				Hand:  ct.Cards(ct.Bind(&nf, Nightforge)),
			},
		})

		h.P1.Play(nf)

		h.P1.ExpectKeys(0)
		h.P1.ExpectAmber(6)
		h.Expect(nf).At(ct.Discard)
	})
}
