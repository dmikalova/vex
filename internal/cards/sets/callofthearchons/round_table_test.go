package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Round Table
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	Each friendly Knight trait creature gains +1 power and taunt.
func TestRoundTable(t *testing.T) {
	t.Run("gives friendly Knights +1 power", func(t *testing.T) {
		var knight, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					RoundTable,
					ct.Bind(&knight, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4), ct.Traits("Knight"))),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4), ct.Traits("Human"))),
				),
			},
		})

		h.Expect(knight).Power(5)
		h.Expect(other).Power(4)
	})
}
