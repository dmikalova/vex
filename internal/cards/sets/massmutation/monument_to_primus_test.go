package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Monument to Primus
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: If Consul Primus is in your discard pile, move 1 Æmber from a creature to another creature. Otherwise, move 1 Æmber from a friendly creature to another friendly creature.
func TestMonumentToPrimus(t *testing.T) {
	t.Run("moves 1 Æmber between friendly creatures", func(t *testing.T) {
		var giver, taker, bystander ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					MonumentToPrimus,
					ct.Bind(&giver, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
					ct.Bind(&taker, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
					ct.Bind(
						&bystander,
						ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4)),
					),
				),
			},
		})
		h.Game().AddAmberOn(giver.ID(), 2)

		h.P1.UseAction(MonumentToPrimus)
		h.P1.ClickCard(taker)

		h.Expect(giver).AmberOn(1)
		h.Expect(taker).AmberOn(1)
		h.Expect(bystander).AmberOn(0)
	})

	t.Run(
		"moves onto an enemy creature when Consul Primus is in your discard pile",
		func(t *testing.T) {
			var giver, enemy, bystander ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Saurian,
					InPlay: ct.Cards(
						MonumentToPrimus,
						ct.Bind(&giver, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
					),
					Discard: ct.Cards(ConsulPrimus),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
						ct.Bind(
							&bystander,
							ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4)),
						),
					),
				},
			})
			h.Game().AddAmberOn(giver.ID(), 2)

			h.P1.UseAction(MonumentToPrimus)
			h.P1.ClickCard(enemy)

			h.Expect(giver).AmberOn(1)
			h.Expect(enemy).AmberOn(1)
			h.Expect(bystander).AmberOn(0)
		},
	)
}
