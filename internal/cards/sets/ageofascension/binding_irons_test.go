package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Binding Irons
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent gains 3 chains.
func TestBindingIrons(t *testing.T) {
	t.Run("gives the opponent 3 chains", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(BindingIrons),
			},
		})

		h.P1.Play(BindingIrons)

		if got := h.Game().State.Chains[1]; got != 3 {
			t.Errorf("opponent chains = %d, want 3", got)
		}
	})
}
