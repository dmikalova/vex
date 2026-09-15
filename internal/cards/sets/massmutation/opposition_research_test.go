package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Opposition Research
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent cannot use creatures to reap during their next turn.
//	Enhance Damage.
func TestOppositionResearch(t *testing.T) {
	t.Run("bars the opponent from reaping next turn, fighting stays open", func(t *testing.T) {
		var foe, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				Hand:   ct.Cards(OppositionResearch),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.Power(3)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(OppositionResearch)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P2.ExpectCannotUseTo(foe, engine.ReapUse)
		h.P2.Fight(foe, ally) // fighting is still allowed
	})
}
