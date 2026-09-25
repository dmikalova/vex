package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// The Shadowsmith
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant • Thief
//
//	Each Mutant creature gains elusive.
func TestTheShadowsmith(t *testing.T) {
	t.Run("grants elusive to each Mutant creature", func(t *testing.T) {
		var mutant, nonMutant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					TheShadowsmith,
					ct.Bind(&mutant, ct.Creature(ct.Traits(card.Traits.Mutant))),
					ct.Bind(&nonMutant, ct.Creature()),
				),
			},
		})

		if !h.Game().HasKeyword(mutant.ID(), card.Keyword.Elusive) {
			t.Error("Mutant creature should have elusive")
		}
		if h.Game().HasKeyword(nonMutant.ID(), card.Keyword.Elusive) {
			t.Error("non-Mutant creature should not have elusive")
		}
	})
}
