package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Snarette
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	At the end of your turn, Snarette captures 1 Æmber from your opponent.
//	Action: Move each Æmber on Snarette to the common supply.
func TestSnarette(t *testing.T) {
	t.Run("captures 1 Æmber from your opponent at end of turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Snarette),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.EndTurn()

		h.Expect(Snarette).AmberOn(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("action moves each captured Æmber to the supply", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Snarette),
			},
			P2: ct.Side{Amber: 3},
		})
		h.P1.EndTurn()
		h.P2.EndTurn()

		h.P1.UseAction(Snarette)

		h.Expect(Snarette).AmberOn(0)
		h.P2.ExpectAmber(2)
	})
}
