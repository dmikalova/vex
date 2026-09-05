package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gargantes Scrapper
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Æmber:  1
//	Traits: Giant
//
//	Alpha.
//	Play: For each Æmber in your pool, deal 3 damage to an enemy creature.
func TestGargantesScrapper(t *testing.T) {
	t.Run("deals 3 damage for each Æmber in your pool", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Amber: 3,
				Hand:  ct.Cards(GargantesScrapper),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Play(GargantesScrapper)

		// Playing adds its Æmber bonus (1) to the pool of 3 before the Play
		// ability resolves, so 4 Æmber deal 12 damage.
		h.Expect(foe).Damage(12)
	})

	t.Run("counts its own Æmber bonus with an empty pool", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Amber: 0,
				Hand:  ct.Cards(GargantesScrapper),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Play(GargantesScrapper)

		// The Æmber bonus (1) is in the pool when the Play ability resolves.
		h.Expect(foe).Damage(3)
	})
}
