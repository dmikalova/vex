package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// City Gates
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: A friendly creature captures 1 Æmber from your opponent. If it is a Dinosaur creature, it captures 1 Æmber from your opponent.
func TestCityGates(t *testing.T) {
	t.Run("a Dinosaur captures 2", func(t *testing.T) {
		var dino ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					CityGates,
					ct.Bind(&dino, ct.Creature(
						ct.OfHouse(card.House.Saurian),
						ct.Traits(card.Traits.Dinosaur),
					)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(CityGates)

		h.Expect(dino).AmberOn(2)
		h.P2.ExpectAmber(3)
	})

	t.Run("a non-Dinosaur captures 1", func(t *testing.T) {
		var soldier ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					CityGates,
					ct.Bind(&soldier, ct.Creature(ct.OfHouse(card.House.Saurian))),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(CityGates)

		h.Expect(soldier).AmberOn(1)
		h.P2.ExpectAmber(4)
	})
}
