package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Crystal Surge
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Exalt each Mutant creature.
func TestCrystalSurge(t *testing.T) {
	t.Run("exalts each Mutant creature only", func(t *testing.T) {
		var mutant, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(CrystalSurge),
				InPlay: ct.Cards(
					ct.Bind(&mutant, ct.Creature(ct.Traits(card.Traits.Mutant))),
					ct.Bind(&plain, ct.Creature()),
				),
			},
		})

		h.P1.Play(CrystalSurge)

		h.Expect(mutant).AmberOn(1)
		h.Expect(plain).AmberOn(0)
	})
}
