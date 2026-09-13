package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Harvest Time
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a Creature - purge each Creature that shares a trait with it. For each card they controlled that was purged this way, each player gains 1 Æmber.
func TestHarvestTime(t *testing.T) {
	t.Run(
		"purges each creature sharing a trait with the chosen creature and pays each player",
		func(t *testing.T) {
			var demon1, demon2, foe, survivor ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					Hand:  ct.Cards(HarvestTime),
					InPlay: ct.Cards(
						ct.Bind(
							&demon1,
							ct.Creature(ct.OfHouse(card.House.Dis), ct.Traits(card.Traits.Demon)),
						),
						ct.Bind(
							&demon2,
							ct.Creature(ct.OfHouse(card.House.Dis), ct.Traits(card.Traits.Demon)),
						),
						ct.Bind(
							&survivor,
							ct.Creature(ct.OfHouse(card.House.Dis), ct.Traits(card.Traits.Beast)),
						),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(
							&foe,
							ct.Creature(
								ct.OfHouse(card.House.Untamed),
								ct.Traits(card.Traits.Demon),
							),
						),
					),
				},
			})

			h.P1.Play(HarvestTime)
			h.P1.ClickCard(demon1)

			h.Expect(demon1).At(ct.Purge)
			h.Expect(demon2).At(ct.Purge)
			h.Expect(foe).At(ct.Purge)
			h.Expect(survivor).At(ct.PlayArea)
			h.P1.ExpectAmber(2)
			h.P2.ExpectAmber(1)
		},
	)
}
