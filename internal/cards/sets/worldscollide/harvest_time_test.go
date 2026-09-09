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
//	Play: Choose a trait, then purge each card with that trait. Each player gains 1 Æmber for each card they controlled that was purged this way.
func TestHarvestTime(t *testing.T) {
	t.Run("purges every card of the chosen trait and pays each player", func(t *testing.T) {
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
						ct.Creature(ct.OfHouse(card.House.Untamed), ct.Traits(card.Traits.Demon)),
					),
				),
			},
		})

		h.P1.Play(HarvestTime)
		h.P1.ClickOption("Demon")

		h.Expect(demon1).At(ct.Purge)
		h.Expect(demon2).At(ct.Purge)
		h.Expect(foe).At(ct.Purge)
		h.Expect(survivor).At(ct.PlayArea)
		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})
}
