package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Grenade Snib
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Goblin
//
//	Destroyed: Your opponent loses 2 Æmber.
func TestGrenadeSnib(t *testing.T) {
	t.Run("makes the opponent lose 2 Æmber when destroyed", func(t *testing.T) {
		var snib, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(ct.Bind(&snib, GrenadeSnib))},
			P2: ct.Side{
				Amber: 3,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Fight(snib, foe)

		h.Expect(snib).At(ct.Discard)
		h.P2.ExpectAmber(1)
	})
}
