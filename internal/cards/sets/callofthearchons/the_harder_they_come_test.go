package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Harder They Come
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Purge a creature with power 5 or higher.
func TestTheHarderTheyCome(t *testing.T) {
	t.Run("purges a creature with power 5 or higher", func(t *testing.T) {
		var strong, weak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(TheHarderTheyCome)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&strong, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(6))),
					ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Play(TheHarderTheyCome)

		h.Expect(strong).At(ct.Purge)
		h.Expect(weak).At(ct.PlayArea)
	})
}
