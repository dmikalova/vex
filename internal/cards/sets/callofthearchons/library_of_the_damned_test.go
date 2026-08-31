package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Library of the Damned
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Action: Archive a card from your hand.
func TestLibraryOfTheDamned(t *testing.T) {
	t.Run("archives a card from hand via its Action", func(t *testing.T) {
		var spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand: ct.Cards(
					ct.Bind(&spare, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(2))),
				),
				InPlay: ct.Cards(LibraryOfTheDamned),
			},
		})

		h.P1.UseAction(LibraryOfTheDamned)

		h.Expect(spare).At(ct.Archives)
	})
}
