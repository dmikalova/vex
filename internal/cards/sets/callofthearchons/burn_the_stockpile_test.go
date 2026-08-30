package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Burn the Stockpile
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: If your opponent has 7 Æmber or more, your opponent loses 4 Æmber.
func TestBurnTheStockpile(t *testing.T) {
	t.Run("drains 4 Æmber from an opponent at 7 or more", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(BurnTheStockpile)},
			P2: ct.Side{Amber: 7},
		})

		h.P1.Play(BurnTheStockpile)

		h.P2.ExpectAmber(3)
	})

	t.Run("does nothing below 7 Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(BurnTheStockpile)},
			P2: ct.Side{Amber: 6},
		})

		h.P1.Play(BurnTheStockpile)

		h.P2.ExpectAmber(6)
	})
}
