package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Library of Babble
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Action: Draw a card.
func TestLibraryOfBabble(t *testing.T) {
	t.Run("draws a card", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(LibraryOfBabble),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
		})

		h.P1.UseAction(LibraryOfBabble)

		h.Expect(top).At(ct.Hand)
	})
}
