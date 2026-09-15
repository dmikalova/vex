package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Looter Goblin
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Goblin
//
//	Elusive.
//	Reap: For the remainder of the turn, each time an enemy creature is destroyed, gain 1 Æmber.
func TestLooterGoblin(t *testing.T) {
	t.Run("reaping gains 1 Æmber per enemy creature destroyed this turn", func(t *testing.T) {
		var goblin, ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&goblin, LooterGoblin),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Reap(goblin)
		h.P1.ExpectAmber(1)

		h.P1.Fight(ally, foe)
		h.Expect(foe).At(ct.Discard)
		h.P1.ExpectAmber(2)
	})
}
