package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Macis Asp
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Skirmish, Poison.
func TestMacisAsp(t *testing.T) {
	t.Run("is a 3-power creature with Skirmish and Poison", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(MacisAsp)},
		})

		h.P1.Play(MacisAsp)

		h.Expect(MacisAsp).Power(3).At(ct.PlayArea)

		var hasSkirmish, hasPoison bool
		for _, kw := range MacisAsp.Keywords {
			switch kw {
			case card.Keyword.Skirmish:
				hasSkirmish = true
			case card.Keyword.Poison:
				hasPoison = true
			}
		}
		if !hasSkirmish || !hasPoison {
			t.Errorf("Macis Asp keywords = %v, want Skirmish and Poison", MacisAsp.Keywords)
		}
	})
}
