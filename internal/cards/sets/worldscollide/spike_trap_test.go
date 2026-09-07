package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Spike Trap
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Spike Trap -> deal 3 damage to each flank creature.
func TestSpikeTrap(t *testing.T) {
	t.Run("destroys itself and deals 3 damage to each flank creature", func(t *testing.T) {
		var trap, left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&trap, SpikeTrap)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(5))),
				ct.Bind(&mid, ct.Creature(ct.Power(5))),
				ct.Bind(&right, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.UseAction(trap)

		h.Expect(trap).At(ct.Discard)
		h.Expect(left).Damage(3)
		h.Expect(right).Damage(3)
		h.Expect(mid).Damage(0)
	})
}
