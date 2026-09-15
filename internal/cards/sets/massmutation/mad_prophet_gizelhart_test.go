package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mad Prophet Gizelhart
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  3
//	Traits: Leader • Priest
//
//	Action: If Mad Prophet Gizelhart is in the center of your battleline, fully heal each non-Mutant creature. For each creature healed this way, gain 1 Æmber.
func TestMadProphetGizelhart(t *testing.T) {
	t.Run(
		"centered: fully heals each non-Mutant creature and gains 1 per creature healed",
		func(t *testing.T) {
			var flankA, gizelhart, flankB, mutant ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Sanctum,
					InPlay: ct.Cards(
						ct.Bind(&flankA, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
						ct.Bind(&gizelhart, MadProphetGizelhart),
						ct.Bind(&flankB, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&mutant, ct.Creature(ct.Traits(card.Traits.Mutant), ct.Power(4))),
					),
				},
			})
			h.Game().State.Cards[flankA.ID()].Damage = 2
			h.Game().State.Cards[flankB.ID()].Damage = 2
			h.Game().State.Cards[mutant.ID()].Damage = 2

			h.P1.UseAction(gizelhart)

			h.Expect(flankA).Damage(0)
			h.Expect(flankB).Damage(0)
			h.Expect(mutant).Damage(2) // a Mutant creature is not healed
			h.P1.ExpectAmber(2)        // one Æmber per non-Mutant creature healed
		},
	)

	t.Run("off-center: heals nothing and gains no Æmber", func(t *testing.T) {
		var gizelhart, flankA ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&gizelhart, MadProphetGizelhart),
					ct.Bind(&flankA, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
		})
		h.Game().State.Cards[flankA.ID()].Damage = 2

		h.P1.UseAction(gizelhart)

		h.Expect(flankA).Damage(2)
		h.P1.ExpectAmber(0)
	})
}
