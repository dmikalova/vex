package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Positron Bolt
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Choose a flank creature. Deal 3 damage to it, 2 damage to its neighbor, and 1 damage to the neighbor's other neighbor.
func TestPositronBolt(t *testing.T) {
	t.Run("deals 3/2/1 walking inward from the chosen flank creature", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(PositronBolt)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(9))),
				ct.Bind(&mid, ct.Creature(ct.Power(9))),
				ct.Bind(&right, ct.Creature(ct.Power(9))),
			)},
		})

		h.P1.Play(PositronBolt)
		h.P1.ClickCard(left)

		h.Expect(left).Damage(3)
		h.Expect(mid).Damage(2)
		h.Expect(right).Damage(1)
	})
}
