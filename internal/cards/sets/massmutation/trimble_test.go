package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Trimble
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant
//
//	Each Mutant creature gains skirmish.
func TestTrimble(t *testing.T) {
	t.Run("grants skirmish to each Mutant creature", func(t *testing.T) {
		var mutant, nonMutant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					Trimble,
					ct.Bind(&mutant, ct.Creature(ct.Traits(card.Traits.Mutant))),
					ct.Bind(&nonMutant, ct.Creature()),
				),
			},
		})

		if !h.Game().HasKeyword(mutant.ID(), card.Keyword.Skirmish) {
			t.Error("Mutant creature should have skirmish")
		}
		if h.Game().HasKeyword(nonMutant.ID(), card.Keyword.Skirmish) {
			t.Error("non-Mutant creature should not have skirmish")
		}
	})
}
