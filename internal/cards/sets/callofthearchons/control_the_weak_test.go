package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Control the Weak
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a house - your opponent must choose that house as their active house during their next turn.
func TestControlTheWeak(t *testing.T) {
	t.Run("forces the opponent's active house on their next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(ControlTheWeak)},
			P2: ct.Side{House: card.House.Mars},
		})

		h.P1.Play(ControlTheWeak)
		h.P1.ExpectPrompt("Choose a house").Source("Control the Weak")
		h.P1.ClickOption("Mars")

		if got := h.Game().State.ForcedHouseNext[1]; got != card.House.Mars {
			t.Fatalf("armed house = %v, want Mars", got)
		}

		h.P1.EndTurn() // the opponent's turn begins, promoting the forced house

		if err := h.Game().
			ChooseHouse(1, card.House.Sanctum); err != engine.ErrMustChooseForcedHouse {
			t.Errorf("wrong house = %v, want ErrMustChooseForcedHouse", err)
		}
		h.P2.ChooseHouse(card.House.Mars) // the forced house is allowed
	})
}
