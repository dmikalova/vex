package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Mutagenic Serum
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Mutagenic Serum. For the remainder of the turn, you may use friendly Mutant creatures.
func TestMutagenicSerum(t *testing.T) {
	t.Run("destroys itself and grants use of friendly Mutant creatures", func(t *testing.T) {
		var mutant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					MutagenicSerum,
					ct.Bind(&mutant, ct.Creature(
						ct.OfHouse(card.House.Untamed),
						ct.Traits(card.Traits.Mutant),
					)),
				),
			},
		})

		// The off-house Mutant creature cannot reap until the serum frees it.
		h.P1.ExpectCannotUseTo(mutant, engine.ReapUse)

		h.P1.UseAction(MutagenicSerum)

		h.Expect(MutagenicSerum).At(ct.Discard)
		if h.Game().State.MayUseTrait[0] != engine.Mutant {
			t.Error("serum should grant use of friendly Mutant creatures this turn")
		}

		h.P1.Reap(mutant)
		h.P1.ExpectAmber(1)
	})

	t.Run("does not free a non-Mutant off-house creature", func(t *testing.T) {
		var beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					MutagenicSerum,
					ct.Bind(&beast, ct.Creature(
						ct.OfHouse(card.House.Untamed),
						ct.Traits(card.Traits.Beast),
					)),
				),
			},
		})

		h.P1.UseAction(MutagenicSerum)

		h.P1.ExpectCannotUseTo(beast, engine.ReapUse)
	})
}
