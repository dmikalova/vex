package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Dark Wave
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to each non-Mutant creature.
func TestDarkWave(t *testing.T) {
	t.Run("deals 2 damage to each non-Mutant creature", func(t *testing.T) {
		var ally, mutant, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Hand:   ct.Cards(DarkWave),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.Power(5)))),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&mutant, ct.Creature(ct.Power(5), ct.Traits(card.Traits.Mutant))),
				ct.Bind(&plain, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(DarkWave)

		h.Expect(mutant).Damage(0)
		h.Expect(plain).Damage(2)
		h.Expect(ally).Damage(2)
	})
}
