package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Longfused Mines
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Longfused Mines, and deal 3 damage to each enemy Creature that is not on a flank.
func TestLongfusedMines(t *testing.T) {
	t.Run("sacrifices itself and deals 3 to each non-flank enemy", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(LongfusedMines),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
					ct.Bind(&mid, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.UseAction(LongfusedMines)

		h.Expect(mid).Damage(3)  // interior enemy is hit
		h.Expect(left).Damage(0) // flanks are spared
		h.Expect(right).Damage(0)
		h.Expect(LongfusedMines).At(ct.Discard) // sacrificed itself
	})
}
