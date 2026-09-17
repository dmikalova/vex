package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Desire
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Demon • Sin
//
//	Each player's keys cost +4 Æmber.
//	Reap: Forge a key at current cost, reduced by 1 Æmber for each friendly Sin creature -> purge Desire.
func TestDesire(t *testing.T) {
	t.Run("forges below the raised cost and purges itself", func(t *testing.T) {
		var desire ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&desire, Desire),
					ct.Creature(
						ct.OfHouse(card.House.Dis),
						ct.Traits(card.Traits.Sin),
						ct.Power(3),
					),
					ct.Creature(
						ct.OfHouse(card.House.Dis),
						ct.Traits(card.Traits.Sin),
						ct.Power(3),
					),
				),
				Amber: 7,
			},
		})

		// The reap gains 1 Æmber (pool 8). The current cost is base 6 plus Desire's
		// own +4 = 10, discounted by the three friendly Sin creatures to 7.
		h.P1.Reap(desire)

		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(1)
		h.Expect(desire).At(ct.Purge) // a landed forge purges the source
	})
}
