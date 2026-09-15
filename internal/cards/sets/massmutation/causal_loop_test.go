package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Causal Loop
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Archive a card from your hand. Archive Causal Loop.
func TestCausalLoop(t *testing.T) {
	t.Run("archives a card and itself", func(t *testing.T) {
		var other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					CausalLoop,
					ct.Bind(&other, ct.Creature()),
				),
			},
		})

		h.P1.Play(CausalLoop)

		h.Expect(other).At(ct.Archives)
		h.Expect(CausalLoop).At(ct.Archives)
	})
}
