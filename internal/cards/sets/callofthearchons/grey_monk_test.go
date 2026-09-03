package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Grey Monk
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human • Priest
//
//	Each friendly creature gains +1 armor.
//	Reap: Heal 2 damage from a creature.
func TestGreyMonk(t *testing.T) {
	t.Run("gives each friendly creature +1 armor", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					GreyMonk,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Armor(0))),
				),
			},
		})

		h.Expect(ally).Armor(1)
	})

	t.Run("heals 2 damage from a creature when it reaps", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					GreyMonk,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(5))),
				),
			},
		})

		ally.Damaged(2)
		h.P1.Reap(GreyMonk)

		h.Expect(ally).Damage(0)
	})
}
