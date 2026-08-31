package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Duma the Martyr
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	Destroyed: Fully heal each other friendly creature, and draw 2 cards.
func TestDumaTheMartyr(t *testing.T) {
	t.Run(
		"fully heals each other friendly creature and draws 2 when destroyed",
		func(t *testing.T) {
			var duma, ally, d1, d2, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Sanctum,
					InPlay: ct.Cards(
						ct.Bind(&duma, DumaTheMartyr),
						ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
					),
					Deck: ct.Cards(
						ct.Bind(&d1, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(1))),
						ct.Bind(&d2, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(1))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
					),
				},
			})
			ally.Damaged(2)

			h.P1.Fight(duma, enemy)

			h.Expect(duma).At(ct.Discard) // 3 power dies to 3 return damage
			h.Expect(ally).Damage(0)      // fully healed
			h.Expect(d1).At(ct.Hand)      // drew 2
			h.Expect(d2).At(ct.Hand)
		},
	)
}
