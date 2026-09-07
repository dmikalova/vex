package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Berserker Slam
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 4 damage to a flank creature. If this damage destroys that creature, its controller loses 1 Æmber.
func TestBerserkerSlam(t *testing.T) {
	t.Run("destroys a flank creature and its controller loses 1 Æmber", func(t *testing.T) {
		var flank ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(BerserkerSlam)},
			P2: ct.Side{
				Amber:  2,
				InPlay: ct.Cards(ct.Bind(&flank, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Play(BerserkerSlam)

		h.Expect(flank).At(ct.Discard)
		h.P2.ExpectAmber(1) // lost 1 because the damage destroyed its creature
	})

	t.Run("no Æmber lost when the creature survives", func(t *testing.T) {
		var flank ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(BerserkerSlam)},
			P2: ct.Side{
				Amber:  2,
				InPlay: ct.Cards(ct.Bind(&flank, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Play(BerserkerSlam)

		h.Expect(flank).At(ct.PlayArea)
		h.P2.ExpectAmber(2)
	})
}
