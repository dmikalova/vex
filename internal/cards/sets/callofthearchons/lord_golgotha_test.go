package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lord Golgotha
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  2
//	Traits: Knight • Spirit
//
//	Before Fight: Deal 3 damage to each neighbor of the creature Lord Golgotha fights.
func TestLordGolgotha(t *testing.T) {
	t.Run("damages the neighbors of the creature it fights, not that creature", func(t *testing.T) {
		var golgotha, left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&golgotha, LordGolgotha)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(9))),
					ct.Bind(&mid, ct.Creature(ct.Power(9))),
					ct.Bind(&right, ct.Creature(ct.Power(9))),
				),
			},
		})

		h.P1.Fight(golgotha, mid)

		h.Expect(left).Damage(3)
		h.Expect(right).Damage(3)
		h.Expect(mid).Damage(5) // fight damage only
	})
}
