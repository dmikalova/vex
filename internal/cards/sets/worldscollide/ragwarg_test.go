package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ragwarg
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	After a creature reaps, if it is the first time a creature has reaped this turn, deal 2 damage to it.
func TestRagwarg(t *testing.T) {
	t.Run("deals 2 damage to the first creature that reaps this turn", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Ragwarg, ct.Bind(&reaper, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Reap(reaper)

		h.Expect(reaper).Damage(2)
		h.P1.ExpectAmber(1)
	})
}
