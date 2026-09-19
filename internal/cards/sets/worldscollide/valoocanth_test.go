package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Valoocanth
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Aquan
//
//	While the tide is low, Valoocanth cannot be used.
//	Fight/Reap: Exhaust an enemy creature and each of its neighbors.
func TestValoocanth(t *testing.T) {
	mars := ct.OfHouse(card.House.Mars)

	t.Run("reaping exhausts a chosen enemy creature and each of its neighbors", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Valoocanth),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(mars)),
				ct.Bind(&mid, ct.Creature(mars)),
				ct.Bind(&right, ct.Creature(mars)),
			)},
		})

		// The tide is neutral by default, so Valoocanth may be used.
		h.P1.Reap(Valoocanth)
		h.P1.ClickCard(mid)

		h.Expect(left).Exhausted()
		h.Expect(mid).Exhausted()
		h.Expect(right).Exhausted()
	})

	t.Run("fighting exhausts the fought enemy and each of its neighbors", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Valoocanth),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(mars, ct.Power(1))),
				// Enough armor to survive the fight, so it can be the exhaust anchor.
				ct.Bind(&mid, ct.Creature(mars, ct.Power(1), ct.Armor(6))),
				ct.Bind(&right, ct.Creature(mars, ct.Power(1))),
			)},
		})

		h.P1.Fight(Valoocanth, mid)
		h.P1.ClickCard(mid)

		h.Expect(left).Exhausted()
		h.Expect(mid).Exhausted()
		h.Expect(right).Exhausted()
	})
}
