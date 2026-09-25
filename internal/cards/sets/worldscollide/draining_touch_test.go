package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Draining Touch
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy a creature with no Æmber on it.
func TestDrainingTouch(t *testing.T) {
	t.Run("destroys a creature with no Æmber on it", func(t *testing.T) {
		var bare, rich ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(DrainingTouch),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&bare, ct.Creature(ct.Power(3))),
				ct.Bind(&rich, ct.Creature(ct.Power(3))),
			)},
		})
		h.Game().State.Cards[rich.ID()].Amber = 1

		// Only the bare creature qualifies, so the choice resolves to it.
		h.P1.Play(DrainingTouch)

		h.Expect(bare).At(ct.Discard)
		h.Expect(rich).At(ct.PlayArea)
	})
}
