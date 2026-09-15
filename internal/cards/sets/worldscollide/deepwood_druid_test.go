package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Deepwood Druid
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Witch
//
//	Deploy.
//	Play/Reap: Fully heal a neighboring creature.
func TestDeepwoodDruid(t *testing.T) {
	t.Run("reaping fully heals a neighboring creature", func(t *testing.T) {
		var druid, neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&druid, DeepwoodDruid),
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(6))),
				),
			},
		})
		neighbor.Damaged(5)

		h.P1.Reap(druid)

		h.Expect(neighbor).Damage(0)
	})
}
