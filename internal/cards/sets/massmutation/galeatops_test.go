package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Galeatops
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  12
//	Traits: Beast
//
//	Galeatops deals 4 Damage when fighting.
func TestGaleatops(t *testing.T) {
	t.Run("only deals 4 damage when fighting", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, InPlay: ct.Cards(Galeatops)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(12))),
				),
			},
		})

		h.P1.Fight(Galeatops, enemy)

		h.Expect(enemy).Damage(4)
	})
}
