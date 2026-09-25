package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Dark Queen Gloriana
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Mutant
//
//	Play: Put a friendly non-Untamed creature into its owner's hand.
//	Enhance Æmber Æmber.
func TestDarkQueenGloriana(t *testing.T) {
	t.Run("returns a friendly non-Untamed creature to hand", func(t *testing.T) {
		var offHouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(DarkQueenGloriana),
				InPlay: ct.Cards(
					ct.Bind(&offHouse, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Play(DarkQueenGloriana)

		h.Expect(offHouse).At(ct.Hand)
	})
}
