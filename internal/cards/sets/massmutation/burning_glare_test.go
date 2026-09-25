package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Burning Glare
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose one:
//	- Stun an enemy creature
//	- Stun each enemy Mutant creature.
//	Enhance Damage.
func TestBurningGlare(t *testing.T) {
	t.Run("stuns one chosen enemy creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(BurningGlare),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Play(BurningGlare)
		h.P1.ClickOption("an enemy creature") // the single-creature option

		h.Expect(foe).Stunned(true)
	})

	t.Run("stuns each enemy Mutant creature", func(t *testing.T) {
		var mutant1, mutant2, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(BurningGlare),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&mutant1, ct.Creature(ct.Power(3), ct.Traits(card.Traits.Mutant))),
				ct.Bind(&mutant2, ct.Creature(ct.Power(3), ct.Traits(card.Traits.Mutant))),
				ct.Bind(&plain, ct.Creature(ct.Power(3))),
			)},
		})

		h.P1.Play(BurningGlare)
		h.P1.ClickOption("each enemy Mutant")

		h.Expect(mutant1).Stunned(true)
		h.Expect(mutant2).Stunned(true)
		h.Expect(plain).Stunned(false)
	})
}
