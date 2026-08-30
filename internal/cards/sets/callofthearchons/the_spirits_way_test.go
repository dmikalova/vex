package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Spirit's Way
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy each creature with power 3 or higher.
func TestTheSpiritsWay(t *testing.T) {
	t.Run("destroys each creature with power 3 or higher", func(t *testing.T) {
		var strong, weak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(TheSpiritsWay)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&strong, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
					ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Play(TheSpiritsWay)

		h.Expect(strong).At(ct.Discard)
		h.Expect(weak).At(ct.PlayArea)
	})
}
