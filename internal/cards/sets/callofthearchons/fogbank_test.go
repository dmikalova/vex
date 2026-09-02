package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fogbank
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Your opponent cannot use creatures to fight during their next turn.
func TestFogbank(t *testing.T) {
	t.Run("bars the opponent from fighting on their next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(Fogbank)},
		})

		h.P1.Play(Fogbank)

		if !h.Game().State.CannotFightNext[1].Value {
			t.Error("Fogbank should arm the opponent's next turn")
		}
		if h.Game().State.CannotFightNext[0].Value {
			t.Error("Fogbank should not restrict the caster")
		}
	})
}
