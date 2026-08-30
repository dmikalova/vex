package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Booby Trap
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 4 damage to a creature that is not on a flank and 2 damage to each of its neighbors.
func TestBoobyTrap(t *testing.T) {
	t.Run("deals 4 to a non-flank creature and 2 to each of its neighbors", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(BoobyTrap)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(10))),
					ct.Bind(&mid, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(10))),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(10))),
				),
			},
		})

		h.P1.Play(BoobyTrap) // Mid is the only non-flank creature, so it is the target

		h.Expect(mid).Damage(4)
		h.Expect(left).Damage(2)
		h.Expect(right).Damage(2)
	})
}
