package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fangtooth Cavern
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	At the end of your turn, destroy the least powerful Creature.
func TestFangtoothCavern(t *testing.T) {
	t.Run("destroys the least powerful creature at the end of your turn", func(t *testing.T) {
		var weak, strong ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					FangtoothCavern,
					ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(2))),
					ct.Bind(&strong, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(6))),
				),
			},
		})

		h.P1.EndTurn()

		h.Expect(weak).At(ct.Discard)
		h.Expect(strong).At(ct.PlayArea)
	})
}
