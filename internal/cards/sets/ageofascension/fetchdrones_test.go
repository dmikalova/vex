package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Fetchdrones
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Discard the top 2 cards of your deck. For each Logos card discarded this way, a friendly creature captures 2 Æmber from your opponent.
func TestFetchdrones(t *testing.T) {
	t.Run("captures 2 Æmber per Logos card discarded", func(t *testing.T) {
		var drones, captor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&drones, Fetchdrones),
					ct.Bind(&captor, ct.Creature(ct.Power(4))),
				),
				Deck: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Logos)),
					ct.Creature(ct.OfHouse(card.House.Logos)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(drones)

		// Both discarded cards are Logos: the sole friendly creature captures 2 twice.
		h.Expect(captor).AmberOn(4)
		h.P2.ExpectAmber(1)
	})

	t.Run("captures nothing for a non-Logos card", func(t *testing.T) {
		var drones, captor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&drones, Fetchdrones),
					ct.Bind(&captor, ct.Creature(ct.Power(4))),
				),
				Deck: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Mars)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(drones)

		h.Expect(captor).AmberOn(0)
		h.P2.ExpectAmber(5)
	})
}
