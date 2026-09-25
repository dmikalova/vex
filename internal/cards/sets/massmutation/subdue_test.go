package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Subdue
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to a creature and stun it.
func TestSubdue(t *testing.T) {
	t.Run("deals 1 damage to a creature and stuns it", func(t *testing.T) {
		var target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(Subdue),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&target, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
			},
		})

		h.P1.Play(Subdue)

		h.Expect(target).At(ct.PlayArea).Damage(1).Stunned(true)
	})
}
