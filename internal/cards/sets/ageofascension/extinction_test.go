package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Extinction
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature - destroy each creature that shares a trait with it. Gain 1 chain.
func TestExtinction(t *testing.T) {
	t.Run("destroys the chosen creature and every creature sharing a trait", func(t *testing.T) {
		var beastA, beastB, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(Extinction),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&beastA, ct.Creature(ct.Traits(card.Traits.Beast))),
				ct.Bind(&beastB, ct.Creature(ct.Traits(card.Traits.Beast))),
				ct.Bind(&other, ct.Creature(ct.Traits(card.Traits.Knight))),
			)},
		})

		h.P1.Play(Extinction)
		h.P1.ClickCard(beastA)

		h.Expect(beastA).At(ct.Discard)
		h.Expect(beastB).At(ct.Discard)
		h.Expect(other).At(ct.PlayArea)
		if got := h.Game().State.Chains[0]; got != 1 {
			t.Fatalf("chains = %d, want 1", got)
		}
	})
}
