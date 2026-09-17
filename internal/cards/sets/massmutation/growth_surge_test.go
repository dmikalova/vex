package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Growth Surge
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a flank creature. Give it three +1 power counters, its neighbor two +1 power counters, and the neighbor's other neighbor a +1 power counter.
func TestGrowthSurge(t *testing.T) {
	t.Run("gives 3/2/1 counters walking inward from the chosen flank creature", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(GrowthSurge),
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(4))),
					ct.Bind(&mid, ct.Creature(ct.Power(4))),
					ct.Bind(&right, ct.Creature(ct.Power(4))),
				),
			},
		})

		h.P1.Play(GrowthSurge)
		h.P1.ClickCard(left)

		h.Expect(left).Power(7)
		h.Expect(mid).Power(6)
		h.Expect(right).Power(5)
	})

	t.Run("stops at the far flank on a short battleline", func(t *testing.T) {
		var left, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(GrowthSurge),
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(4))),
					ct.Bind(&right, ct.Creature(ct.Power(4))),
				),
			},
		})

		h.P1.Play(GrowthSurge)
		h.P1.ClickCard(left)

		h.Expect(left).Power(7)
		h.Expect(right).Power(6)
	})
}
