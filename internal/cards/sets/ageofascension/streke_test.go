package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Streke
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	While Streke is not on a flank, your opponent's hand size is 1 less.
func TestStreke(t *testing.T) {
	t.Run("slows the opponent's refill while off a flank", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Creature(ct.Power(3)),
					Streke,
					ct.Creature(ct.Power(3)),
				),
			},
			P2: ct.Side{Deck: ct.DeckOf(card.House.Logos, 10)},
		})

		h.P1.EndTurn()
		h.P2.EndTurn()

		if got := int(h.Game().State.Hand[1].Count); got != 5 {
			t.Errorf("opponent hand after draw = %d, want 5", got)
		}
	})

	t.Run("does nothing while on a flank", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Streke),
			},
			P2: ct.Side{Deck: ct.DeckOf(card.House.Logos, 10)},
		})

		h.P1.EndTurn()
		h.P2.EndTurn()

		if got := int(h.Game().State.Hand[1].Count); got != 6 {
			t.Errorf("opponent hand after draw = %d, want 6", got)
		}
	})
}
