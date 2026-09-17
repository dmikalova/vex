package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ultra Gravitron
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  10
//	Armor:  3
//	Traits: Robot
//
//	Play: Archive the top 5 cards of your deck.
//	Fight/Reap: Discard a card from your archives -> purge a creature, and resolve that card's bonus icons.
func TestUltraGravitron(t *testing.T) {
	t.Run("play archives the top 5 cards of the deck", func(t *testing.T) {
		var ug ct.Card
		var deck [5]ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(ct.Bind(&ug, UltraGravitron), card.GiganticArt(UltraGravitron)),
				Deck: ct.Cards(
					ct.Bind(&deck[0], ct.Creature(ct.Power(1))),
					ct.Bind(&deck[1], ct.Creature(ct.Power(1))),
					ct.Bind(&deck[2], ct.Creature(ct.Power(1))),
					ct.Bind(&deck[3], ct.Creature(ct.Power(1))),
					ct.Bind(&deck[4], ct.Creature(ct.Power(1))),
				),
			},
		})

		h.P1.Play(ug)

		for _, c := range deck {
			h.Expect(c).At(ct.Archives)
		}
	})

	t.Run(
		"reap discards from archives, purges a creature, and resolves that card's bonus icons",
		func(t *testing.T) {
			var ug, archived, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Logos,
					InPlay: ct.Cards(ct.Bind(&ug, UltraGravitron)),
					Archives: ct.Cards(
						ct.Bind(&archived, ct.Creature(ct.AemberBonus(1), ct.Power(2))),
					),
				},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				)},
			})

			h.P1.Reap(ug)
			h.P1.ClickCard(enemy) // purge the enemy creature

			h.Expect(archived).At(ct.Discard)
			h.Expect(enemy).At(ct.Purge)
			// The discarded card's Æmber pip resolves as if played, gaining 1.
			h.P1.ExpectAmber(1)
		},
	)
}
