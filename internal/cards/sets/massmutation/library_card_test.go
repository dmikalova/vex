package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Library Card
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: For the remainder of the turn, each time you play another card, draw a card. Purge Library Card.
func TestLibraryCard(t *testing.T) {
	t.Run("purges itself and draws a card each time another card is played", func(t *testing.T) {
		var c1, d1 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(LibraryCard),
				Hand: ct.Cards(
					ct.Bind(&c1, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
				Deck: ct.Cards(
					ct.Bind(&d1, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		// Using the action purges Library Card but draws nothing yet.
		h.P1.UseAction(LibraryCard)
		h.Expect(LibraryCard).At(ct.Purge)
		h.Expect(d1).At(ct.Deck)

		// Playing another card this turn draws.
		h.P1.Play(c1)
		h.Expect(d1).At(ct.Hand)
	})
}
