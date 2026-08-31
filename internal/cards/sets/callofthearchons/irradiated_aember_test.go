package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Irradiated Aember
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If your opponent has 6 Æmber or more, deal 3 damage to each enemy creature.
func TestIrradiatedAember(t *testing.T) {
	t.Run("deals 3 to each enemy creature when the opponent has 6 or more", func(t *testing.T) {
		var toughFoe, weakFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(IrradiatedAember)},
			P2: ct.Side{
				Amber: 6,
				InPlay: ct.Cards(
					ct.Bind(&toughFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
					ct.Bind(&weakFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
			},
		})

		h.P1.Play(IrradiatedAember)

		h.Expect(weakFoe).At(ct.Discard)             // 2 power dies
		h.Expect(toughFoe).At(ct.PlayArea).Damage(3) // 5 power survives marked
	})

	t.Run("does nothing when the opponent is below 6", func(t *testing.T) {
		var survivor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(IrradiatedAember)},
			P2: ct.Side{
				Amber: 5,
				InPlay: ct.Cards(
					ct.Bind(&survivor, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
			},
		})

		h.P1.Play(IrradiatedAember)

		h.Expect(survivor).At(ct.PlayArea).Damage(0)
	})
}
