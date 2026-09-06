package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Angwish
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Demon
//
//	Your opponent's keys cost +1 Æmber for each damage on it.
func TestAngwish(t *testing.T) {
	t.Run("charges the opponent 1 more per damage on Angwish", func(t *testing.T) {
		var angwish ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&angwish, Angwish)),
			},
		})

		if got := h.Game().CurrentKeyCost(1); got != 6 {
			t.Errorf("key cost with no damage = %d, want 6", got)
		}

		angwish.Damaged(2)
		if got := h.Game().CurrentKeyCost(1); got != 8 {
			t.Errorf("key cost with 2 damage = %d, want 8", got)
		}
		if got := h.Game().CurrentKeyCost(0); got != 6 {
			t.Errorf("the controller's own key cost = %d, want 6", got)
		}
	})
}
