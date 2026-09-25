package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// We Can ALL Win
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Each player's keys cost -2 Æmber until the end of your next turn.
func TestWeCanALLWin(t *testing.T) {
	t.Run("drops both players' key cost by 2", func(t *testing.T) {
		var win ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ct.Bind(&win, WeCanALLWin)),
			},
		})

		h.P1.Play(win)

		base := 6
		if got := h.Game().CurrentKeyCost(0); got != base-2 {
			t.Errorf("controller key cost = %d, want %d", got, base-2)
		}
		if got := h.Game().CurrentKeyCost(1); got != base-2 {
			t.Errorf("opponent key cost = %d, want %d", got, base-2)
		}
	})
}
