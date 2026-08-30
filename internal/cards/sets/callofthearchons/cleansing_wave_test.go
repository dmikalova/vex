package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cleansing Wave
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Heal 1 damage from each creature, and for each creature healed this way, gain 1 Æmber.
func TestCleansingWave(t *testing.T) {
	t.Run("heals 1 from each creature and gains 1 Æmber per creature healed", func(t *testing.T) {
		var a, b, healthy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(CleansingWave),
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(5))),
					ct.Bind(&healthy, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(5)))),
			},
		})
		a.Damaged(3)
		b.Damaged(1)

		h.P1.Play(CleansingWave)

		h.Expect(a).Damage(2)
		h.Expect(b).Damage(0)
		h.Expect(healthy).Damage(0) // an undamaged creature stays at 0 and is not counted
		h.P1.ExpectAmber(2)         // 2 creatures were actually healed
	})
}
