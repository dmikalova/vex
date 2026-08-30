package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Save the Pack
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each damaged creature. Gain 1 chain.
func TestSaveThePack(t *testing.T) {
	t.Run("destroys each damaged creature and gains 1 chain", func(t *testing.T) {
		var hurt, healthy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(SaveThePack)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&hurt, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
					ct.Bind(&healthy, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		hurt.Damaged(1)
		h.P1.Play(SaveThePack)

		h.Expect(hurt).At(ct.Discard)
		h.Expect(healthy).At(ct.PlayArea)
		if got := h.Game().State.Chains[0]; got != 1 {
			t.Fatalf("chains = %d, want 1", got)
		}
	})
}
