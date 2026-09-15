package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Crash Muldoon
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Pilot
//
//	Deploy.
//	Crash Muldoon enters play ready.
//	Action: Use a neighboring non-Star Alliance creature.
func TestCrashMuldoon(t *testing.T) {
	t.Run("uses a neighboring non-Star Alliance creature", func(t *testing.T) {
		var crash, neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Bind(&crash, CrashMuldoon),
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(crash)

		h.Expect(neighbor).Exhausted()
		h.P1.ExpectAmber(1) // the neighbor reaps
	})
}
