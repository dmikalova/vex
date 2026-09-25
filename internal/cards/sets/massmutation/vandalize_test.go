package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Vandalize
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Reveal the top 3 cards of your opponent's deck. Discard 1. Put them back in any order.
func TestVandalize(t *testing.T) {
	t.Run("discards 1 of the opponent's top 3 cards, the rest stay", func(t *testing.T) {
		var top, second, third ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Vandalize),
			},
			P2: ct.Side{Deck: ct.Cards(
				ct.Bind(&top, ct.Creature(ct.Power(3))),
				ct.Bind(&second, ct.Creature(ct.Power(4))),
				ct.Bind(&third, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(Vandalize)
		h.P1.ClickCard(top)    // discard one of the three
		h.P1.ClickCard(second) // order the two survivors back onto the deck

		h.Expect(top).At(ct.Discard)
		if deck := h.Game().Deck(1); len(deck) != 2 {
			t.Errorf("opponent deck = %v, want 2 cards remaining", deck)
		}
	})
}
