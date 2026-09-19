package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Particle Sweep
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a creature. If it is a Mutant creature, destroy the chosen creature. Otherwise, deal 2 damage to the chosen creature.
func TestParticleSweep(t *testing.T) {
	t.Run("destroys a Mutant creature instead of damaging it", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ParticleSweep),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.Power(4), ct.Traits(card.Traits.Mutant))),
			)},
		})

		h.P1.Play(ParticleSweep)

		h.Expect(foe).At(ct.Discard)
	})

	t.Run("deals 2 damage to a non-Mutant creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ParticleSweep),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.Power(4))),
			)},
		})

		h.P1.Play(ParticleSweep)

		h.Expect(foe).At(ct.PlayArea).Damage(2)
	})
}
