package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lieutenant Gorvenal
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Spirit • Knight
//
//	After a friendly creature is used to fight, Lieutenant Gorvenal captures 1 Æmber from your opponent.
func TestLieutenantGorvenal(t *testing.T) {
	t.Run("captures 1 Æmber after another friendly creature fights", func(t *testing.T) {
		var ally, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					LieutenantGorvenal,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(1))),
				),
				Amber: 2,
			},
		})

		h.P1.Fight(ally, enemy)

		h.Expect(LieutenantGorvenal).AmberOn(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("does not capture after an enemy creature fights", func(t *testing.T) {
		var gorvenal, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&gorvenal, LieutenantGorvenal),
					ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(1)),
				),
				Amber: 2,
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
			},
		})
		h.P1.EndTurn()

		h.P2.Fight(enemy, gorvenal)

		h.Expect(gorvenal).AmberOn(0)
	})
}
