package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Data Forge
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Forge a key at +10 Æmber current cost, reduced by 1 Æmber for each card in your hand.
func TestDataForge(t *testing.T) {
	t.Run("forges at +10 reduced by 1 per card in hand", func(t *testing.T) {
		var df ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Amber: 20,
				Hand: ct.Cards(
					ct.Bind(&df, DataForge),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
				),
			},
		})

		// After playing Data Forge, 3 cards remain in hand, so +10 reduced by 3
		// is a +7 surcharge on the 6 key cost = 13, paid from 20 + its own 1 Æmber.
		h.P1.Play(df)

		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(8)
	})

	t.Run("does not forge when the surcharge is unaffordable", func(t *testing.T) {
		var df ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Amber: 2,
				Hand:  ct.Cards(ct.Bind(&df, DataForge)),
			},
		})

		// No cards remain, so the surcharge stays +10 on the 6 key cost = 16, more
		// than the 2 + its own 1 Æmber can pay, so no key is forged.
		h.P1.Play(df)

		h.P1.ExpectKeys(0)
		h.P1.ExpectAmber(3)
	})
}
