package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Kaloch Stonefather
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant • Leader
//
//	While Kaloch Stonefather is in the center of your battleline, each friendly creature gains skirmish.
func TestKalochStonefather(t *testing.T) {
	t.Run("centered, friendly creatures gain skirmish and take no retaliation", func(t *testing.T) {
		var fighter, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(
				ct.Bind(&fighter, ct.Creature(ct.Power(10))),
				KalochStonefather,
				ct.Creature(),
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(fighter, enemy)

		h.Expect(fighter).Damage(0)
	})

	t.Run("off-center, friendly creatures take retaliation damage", func(t *testing.T) {
		var fighter, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(
				ct.Bind(&fighter, ct.Creature(ct.Power(10))),
				KalochStonefather,
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(fighter, enemy)

		h.Expect(fighter).Damage(3)
	})
}
