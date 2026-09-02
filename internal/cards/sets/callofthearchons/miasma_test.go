package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Miasma
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Your opponent skips the "forge a key" step during their next turn.
func TestMiasma(t *testing.T) {
	t.Run("makes the opponent skip their next forge step", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(Miasma)},
		})

		h.P1.Play(Miasma)

		if !h.Game().State.SkipForgeNext[1].Value {
			t.Error("the opponent should be set to skip their next forge step")
		}
	})
}
