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

		h.Expect(foe).Damage(9)
	})

	t.Run("deals no damage with an empty pool", func(t *testing.T) {
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

		h.Expect(foe).Damage(0)
	})
}
