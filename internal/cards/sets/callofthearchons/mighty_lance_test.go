package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mighty Lance
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Deal 3 damage to a creature and 3 damage to a neighbor of that creature.
func TestMightyLance(t *testing.T) {
	t.Run("deals 3 to a creature and 3 to a chosen neighbor", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(MightyLance)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(10))),
				ct.Bind(&mid, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(10))),
				ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(10))),
			)},
		})

		h.P1.Play(MightyLance)
		h.P1.ClickCard(mid)  // the creature
		h.P1.ClickCard(left) // the neighbor

		h.Expect(mid).Damage(3)
		h.Expect(left).Damage(3)
		h.Expect(right).Damage(0)
	})
}
