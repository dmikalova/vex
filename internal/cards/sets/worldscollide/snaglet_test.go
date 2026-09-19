package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Snaglet
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Action: Choose a house. If your opponent chooses that house as their active house during their next turn, steal 2 Æmber.
func TestSnaglet(t *testing.T) {
	t.Run("opponent picks the predicted house", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Snaglet),
			},
			P2: ct.Side{
				House: card.House.Mars,
				Amber: 3,
			},
		})

		h.P1.UseAction(Snaglet)
		h.P1.ClickOption("Mars")
		if got := h.Game().State.HouseConstraintsNext[1]; h.Game().State.HouseConstraintCountNext[1] != 1 ||
			got[0].House != card.House.Mars ||
			got[0].Amount != 2 {
			t.Fatalf("armed wager = %+v, want Mars for 2", got[0])
		}

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Mars) // the predicted house pays off

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})

	t.Run("opponent picks a different house", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Snaglet),
			},
			P2: ct.Side{
				House: card.House.Mars,
				Amber: 3,
			},
		})

		h.P1.UseAction(Snaglet)
		h.P1.ClickOption("Mars")

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Logos) // not the predicted house

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
	})
}
