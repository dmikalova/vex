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
//	Play: For each Æmber in your pool, give a creature a +1 power counter.
func TestRapidEvolution(t *testing.T) {
	t.Run("adds a power counter for each Æmber the controller has", func(t *testing.T) {
		var target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(RapidEvolution),
				InPlay: ct.Cards(ct.Bind(&target, ct.Creature(ct.Power(4)))),
				Amber:  3,
			},
		})

		h.P1.Play(RapidEvolution)

		h.Expect(target).Power(8)
	})
}
