package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Graphton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Robot • Mutant
//
//	Reap: Archive the top card of your deck.
func TestGraphton(t *testing.T) {
	t.Run("archives the top card of your deck when it reaps", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Graphton),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Reap(Graphton)

		h.Expect(top).At(ct.Archives)
	})
}
