package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Skippy Timehog
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Play: Your opponent cannot use any cards during their next turn.
func TestSkippyTimehog(t *testing.T) {
	t.Run("stops the opponent reaping on their next turn", func(t *testing.T) {
		var skippy, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(ct.Bind(&skippy, SkippyTimehog)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(skippy)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P2.ExpectCannotUse(enemy)
	})

	t.Run("lifts once that turn is over", func(t *testing.T) {
		var skippy, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(ct.Bind(&skippy, SkippyTimehog)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(skippy)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P2.Reap(enemy)
	})
}
