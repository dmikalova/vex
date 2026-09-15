package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// King of the Crag
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Giant
//
//	Each enemy Brobnar creature gains -2 power.
func TestKingOfTheCrag(t *testing.T) {
	t.Run("gives each enemy Brobnar creature -2 power", func(t *testing.T) {
		var brobFoe, marsFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(KingOfTheCrag)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&brobFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
					ct.Bind(&marsFoe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.Expect(brobFoe).Power(3)
		h.Expect(marsFoe).Power(5)
	})
}
