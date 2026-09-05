package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Martian Hounds
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: For each damaged creature in play, give a creature 2 +1 power counters.
func TestMartianHounds(t *testing.T) {
	t.Run("gives two counters per damaged creature", func(t *testing.T) {
		var hounds, chosen, damaged, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&hounds, MartianHounds)),
				InPlay: ct.Cards(
					ct.Bind(&chosen, ct.Creature(ct.Power(3))),
					ct.Bind(&damaged, ct.Creature(ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(5)))),
			},
		})
		damaged.Damaged(1)
		enemy.Damaged(2)

		h.P1.Play(hounds)
		h.P1.ClickCard(chosen)

		// Two damaged creatures, two +1 counters each.
		h.Expect(chosen).Power(7)
	})

	t.Run("gives nothing when no creature is damaged", func(t *testing.T) {
		var hounds, chosen ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				Hand:   ct.Cards(ct.Bind(&hounds, MartianHounds)),
				InPlay: ct.Cards(ct.Bind(&chosen, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Play(hounds)

		h.Expect(chosen).Power(3)
	})
}
