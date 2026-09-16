package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Krrrzzzaaap!!!
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy each non-Mutant creature, and gain 1 chain.
func TestKrrrzzzaaap(t *testing.T) {
	t.Run("destroys each non-Mutant creature and gains a chain", func(t *testing.T) {
		var mutant, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(Krrrzzzaaap)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&mutant, ct.Creature(ct.Power(3), ct.Traits(card.Traits.Mutant))),
				ct.Bind(&plain, ct.Creature(ct.Power(3))),
			)},
		})

		h.P1.Play(Krrrzzzaaap)

		h.Expect(mutant).At(ct.PlayArea)
		h.Expect(plain).At(ct.Discard)
		if got := h.Game().State.Chains[0]; got != 1 {
			t.Errorf("caster chains = %d, want 1", got)
		}
	})
}
