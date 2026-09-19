package callofthearchons

import (
	"errors"
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
//	Bonus:  Æmber
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

		if got := h.Game().State.HouseConstraintsNext[1]; h.Game().State.HouseConstraintCountNext[1] != 1 ||
			got[0].House != card.House.Mars {
			t.Fatalf("armed constraint = %+v (count %d), want one on Mars",
				got[0], h.Game().State.HouseConstraintCountNext[1])
		}

		h.P1.EndTurn() // the opponent's turn begins, promoting the forced house

		if err := h.Game().
			ChooseHouse(1, card.House.Sanctum); !errors.Is(err, engine.ErrHouseNotAllowed) {
			t.Errorf("wrong house = %v, want ErrHouseNotAllowed", err)
		}
		h.P2.ChooseHouse(card.House.Mars) // the forced house is allowed
	})
}
