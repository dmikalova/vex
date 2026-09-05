package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Scrambler Storm
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Your opponent cannot play Tactics during their next turn.
func TestScramblerStorm(t *testing.T) {
	t.Run("bars the opponent from playing Tactics next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(ScramblerStorm)},
		})

		h.P1.Play(ScramblerStorm)

		if got := h.Game().State.CannotPlayTypeNext[1].Value; got != engine.Tactic {
			t.Errorf("opponent's armed play bar = %q, want Tactic", got)
		}
	})
}
