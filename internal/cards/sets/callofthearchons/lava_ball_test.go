package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lava Ball
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Deal 4 damage to a creature that is not on a flank and 2 damage to each of its neighbors.
func TestLavaBall(t *testing.T) {
	t.Run("deals 4 to a non-flank creature and 2 to its neighbors", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(LavaBall)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(6))),
				ct.Bind(&mid, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(6))),
				ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(6))),
			)},
		})

		h.P1.Play(LavaBall)

		h.Expect(mid).Damage(4)
		h.Expect(left).Damage(2)
		h.Expect(right).Damage(2)
	})
}
