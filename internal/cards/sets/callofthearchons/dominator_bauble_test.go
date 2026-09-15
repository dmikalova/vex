package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dominator Bauble
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	Action: Use a friendly creature.
func TestDominatorBauble(t *testing.T) {
	t.Run("uses the only friendly creature to reap for 1 Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					DominatorBauble,
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
				),
			},
		})

		h.P1.UseAction(DominatorBauble)

		h.P1.ExpectAmber(1)
	})

	// A creature can only be used while ready, so with every friendly creature
	// exhausted there is nothing to choose and nothing is asked.
	t.Run("asks nothing when no friendly creature is ready", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					DominatorBauble,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		ally.Exhaust()
		h.P1.UseAction(DominatorBauble)

		h.P1.ExpectAmber(0)
	})
}
