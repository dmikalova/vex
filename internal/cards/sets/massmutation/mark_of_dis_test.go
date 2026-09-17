package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Mark of Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to a creature. If it is not destroyed, its controller must choose that creature's house as their active house on their next turn.
func TestMarkOfDis(t *testing.T) {
	t.Run("survivor's controller is locked to its house next turn", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(MarkOfDis),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(5))),
				),
			},
		})

		h.P1.Play(MarkOfDis)

		h.Expect(foe).At(ct.PlayArea).Damage(2)
		if h.Game().State.HouseConstraintCountNext[1] != 1 {
			t.Fatalf("armed constraints = %d, want 1", h.Game().State.HouseConstraintCountNext[1])
		}

		h.P1.EndTurn() // the opponent's turn begins, promoting the must

		// The must reads the creature's house live (Logos), so any other house is
		// rejected and Logos is required.
		if err := h.Game().ChooseHouse(1, card.House.Mars); err != engine.ErrHouseNotAllowed {
			t.Errorf("a house other than the creature's = %v, want ErrHouseNotAllowed", err)
		}
		h.P2.ChooseHouse(card.House.Logos)
	})

	t.Run("a destroyed creature arms no constraint", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(MarkOfDis),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
				),
			},
		})

		h.P1.Play(MarkOfDis)

		h.Expect(foe).At(ct.Discard)
		if h.Game().State.HouseConstraintCountNext[1] != 0 {
			t.Errorf("armed constraints = %d, want 0", h.Game().State.HouseConstraintCountNext[1])
		}
	})
}
