package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Stealth Mode
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Until the end of your next turn, players cannot play tactics.
func TestStealthMode(t *testing.T) {
	t.Run(
		"bars both players from playing Tactics until the end of the caster's next turn",
		func(t *testing.T) {
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.StarAlliance,
					Hand:  ct.Cards(StealthMode),
				},
			})

			h.P1.Play(StealthMode)

			if got := h.Game().State.CannotPlayTypeThis[0].Value; got != engine.Tactic {
				t.Errorf("caster's play bar this turn = %q, want Tactic", got)
			}
			if got := h.Game().State.CannotPlayTypeNext[0].Value; got != engine.Tactic {
				t.Errorf("caster's armed play bar = %q, want Tactic", got)
			}
			if got := h.Game().State.CannotPlayTypeNext[1].Value; got != engine.Tactic {
				t.Errorf("opponent's armed play bar = %q, want Tactic", got)
			}
		},
	)
}
