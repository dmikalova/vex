package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Red Alert
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For each creature your opponent controls in excess of you, deal 1 damage to each enemy creature.
func TestRedAlert(t *testing.T) {
	t.Run("deals damage equal to the enemy creature surplus", func(t *testing.T) {
		var strong, weak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(RedAlert),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&strong, ct.Creature(ct.Power(6))),
				ct.Bind(&weak, ct.Creature(ct.Power(4))),
			)},
		})

		// P1 controls no creatures, P2 controls two: a surplus of 2.
		h.P1.Play(RedAlert)

		h.Expect(strong).Damage(2)
		h.Expect(weak).Damage(2)
	})

	t.Run("does nothing when the opponent has no surplus", func(t *testing.T) {
		var friend, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				Hand:   ct.Cards(RedAlert),
				InPlay: ct.Cards(ct.Bind(&friend, ct.Creature(ct.Power(4)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4))))},
		})

		h.P1.Play(RedAlert)

		h.Expect(foe).Damage(0)
	})
}
