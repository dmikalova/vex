package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Impspector
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Destroyed: Purge a random card from your opponent's hand.
func TestImpspector(t *testing.T) {
	t.Run("purges a random card from the opponent's hand when destroyed", func(t *testing.T) {
		var impspector, enemy, doomed ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&impspector, Impspector)),
			},
			P2: ct.Side{
				Hand: ct.Cards(ct.Bind(&doomed, ct.Creature(ct.Power(3)))),
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.Power(2))),
				),
			},
		})

		h.P1.Fight(impspector, enemy)

		h.Expect(impspector).At(ct.Discard) // 2 power dies to 2 return damage
		h.Expect(doomed).At(ct.Purge)       // the sole hand card is purged
	})
}
