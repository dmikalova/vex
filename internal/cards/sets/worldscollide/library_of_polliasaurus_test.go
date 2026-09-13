package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Library of Polliasaurus
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Action: Move 1 Æmber from a friendly Creature to your pool.
func TestLibraryOfPolliasaurus(t *testing.T) {
	t.Run("moves 1 Æmber from a friendly creature to your pool", func(t *testing.T) {
		var lib, banker ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&lib, LibraryOfPolliasaurus),
					ct.Bind(&banker, ct.Creature(ct.OfHouse(card.House.Saurian))),
				),
			},
		})
		h.Game().AddAmberOn(banker.ID(), 2)

		h.P1.UseAction(lib)

		h.Expect(banker).AmberOn(1)
		h.P1.ExpectAmber(1)
	})
}
