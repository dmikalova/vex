package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hallowed Blaster
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Heal 3 damage from a creature.
func TestHallowedBlaster(t *testing.T) {
	t.Run("heals 3 damage from a chosen creature", func(t *testing.T) {
		var wounded ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					HallowedBlaster,
					ct.Bind(&wounded, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(6))),
				),
			},
		})
		wounded.Damaged(4)

		h.P1.UseAction(HallowedBlaster)

		h.Expect(wounded).Damage(1) // 4 - 3
	})
}
