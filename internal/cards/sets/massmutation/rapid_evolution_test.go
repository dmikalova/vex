package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Rapid Evolution
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose a creature - for each Æmber in your pool, give the chosen creature a +1 power counter.
func TestRapidEvolution(t *testing.T) {
	t.Run("stacks every counter on the one chosen creature", func(t *testing.T) {
		var target, bystander ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(RapidEvolution),
				InPlay: ct.Cards(
					ct.Bind(&target, ct.Creature(ct.Power(4))),
					ct.Bind(&bystander, ct.Creature(ct.Power(4))),
				),
				Amber: 3,
			},
		})

		h.P1.Play(RapidEvolution)
		h.P1.ClickCard(target)

		// The Æmber bonus resolves first, so the pool is 4 when the counters land.
		h.Expect(target).Power(8)
		h.Expect(bystander).Power(4)
	})
}
