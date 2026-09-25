package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Daughter
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	Elusive.
//	Your hand size is 1 more.
func TestDaughter(t *testing.T) {
	t.Run("refills its controller's hand to one additional card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Daughter),
				Deck:   ct.DeckOf(card.House.Logos, 10),
			},
		})

		h.P1.EndTurn()

		if got := int(h.Game().State.Hand[0].Count); got != 7 {
			t.Errorf("hand after draw = %d, want 7", got)
		}
	})
}
