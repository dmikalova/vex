package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Stealth Mode
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Your opponent cannot play Tactics during their next turn.
func TestStealthMode(t *testing.T) {
	t.Run("bars the opponent from playing Tactics next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.StarAlliance, Hand: ct.Cards(StealthMode)},
		})

		h.P1.Play(StealthMode)

		if got := h.Game().State.CannotPlayTypeNext[1].Value; got != engine.Tactic {
			t.Errorf("opponent's armed play bar = %q, want Tactic", got)
		}
	})
}
