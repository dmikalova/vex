package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Greed
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Demon • Sin
//
//	For each friendly Sin creature your hand size is 1 more.
func TestGreed(t *testing.T) {
	t.Run("refills one extra card for each friendly Sin creature", func(t *testing.T) {
		var deck []any
		for range 12 {
			deck = append(deck, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(2)))
		}
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					Greed,
					ct.Creature(
						ct.OfHouse(card.House.Dis),
						ct.Traits(card.Traits.Sin),
						ct.Power(3),
					),
				),
				Deck: ct.Cards(deck...),
			},
		})

		h.P1.EndTurn()

		// Two friendly Sin creatures raise the six-card refill to eight.
		if got := len(h.Game().Hand(0)); got != 8 {
			t.Errorf("hand = %d, want 8 (6 + 2 Sin creatures)", got)
		}
	})
}
