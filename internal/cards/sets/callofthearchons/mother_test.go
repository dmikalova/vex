package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mother
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Robot • Scientist
//
//	During your "draw cards" step, refill your hand to 1 additional card.
func TestMother(t *testing.T) {
	t.Run("refills its controller's hand to one additional card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Mother),
				Deck:   ct.DeckOf(card.House.Logos, 10),
			},
		})

		h.P1.EndTurn()

		if got := int(h.Game().State.Hand[0].Count); got != 7 {
			t.Errorf("hand after draw = %d, want 7", got)
		}
	})
}
