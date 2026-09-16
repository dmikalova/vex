package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pismire
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	If you control more Mutant creatures than your opponent, your opponent's keys cost +2 Æmber.
func TestPismire(t *testing.T) {
	t.Run("taxes the opponent while you control more Mutant creatures", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(Pismire)},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(ct.Power(2)))},
		})

		if got := h.Game().CurrentKeyCost(1); got != 6+2 {
			t.Errorf("opponent key cost = %d, want 8", got)
		}
	})

	t.Run("does not tax once Mutant counts are level", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(Pismire)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Creature(ct.Traits(card.Traits.Mutant))),
			},
		})

		if got := h.Game().CurrentKeyCost(1); got != 6 {
			t.Errorf("opponent key cost = %d, want 6", got)
		}
	})
}
