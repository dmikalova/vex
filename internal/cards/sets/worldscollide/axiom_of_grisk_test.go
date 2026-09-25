package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Axiom of Grisk
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ward a creature. Destroy each creature with no Æmber on it. Gain 2 chains.
func TestAxiomOfGrisk(t *testing.T) {
	t.Run("wards one, destroys the rest without Æmber, gains 2 chains", func(t *testing.T) {
		var warded, bare, rich ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(AxiomOfGrisk),
				InPlay: ct.Cards(
					ct.Bind(&warded, ct.Creature(ct.Power(3))),
					ct.Bind(&rich, ct.Creature(ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&bare, ct.Creature(ct.Power(3))))},
		})
		h.Game().State.Cards[rich.ID()].Amber = 1

		h.P1.Play(AxiomOfGrisk)
		h.P1.ClickCard(warded) // ward this one so the destroy spares it

		h.Expect(warded).At(ct.PlayArea) // ward absorbed the destroy
		h.Expect(rich).At(ct.PlayArea)   // has Æmber, not targeted
		h.Expect(bare).At(ct.Discard)    // no Æmber, destroyed
		if got := h.Game().State.Chains[0]; got != 2 {
			t.Errorf("P1 chains = %d, want 2", got)
		}
	})
}
