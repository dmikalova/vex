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
//	Play: Archive 2 cards from your hand. Archive Causal Loop.
func TestCausalLoop(t *testing.T) {
	t.Run("archives two cards and itself", func(t *testing.T) {
		var first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					CausalLoop,
					ct.Bind(&first, ct.Creature()),
					ct.Bind(&second, ct.Creature()),
				),
			},
		})

		h.P1.Play(CausalLoop)
		// The second archive has only one card left, so it needs no click.
		h.P1.ClickCard(first)

		h.Expect(first).At(ct.Archives)
		h.Expect(second).At(ct.Archives)
		h.Expect(CausalLoop).At(ct.Archives)
	})
}
