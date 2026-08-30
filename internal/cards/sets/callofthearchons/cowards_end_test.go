package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Coward's End
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each undamaged creature. Gain 3 chains.
func TestCowardsEnd(t *testing.T) {
	t.Run("destroys each undamaged creature and gains 3 chains", func(t *testing.T) {
		var healthy, hurt ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(CowardsEnd)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&healthy, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&hurt, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		hurt.Damaged(1)
		h.P1.Play(CowardsEnd)

		h.Expect(healthy).At(ct.Discard)
		h.Expect(hurt).At(ct.PlayArea)
		if got := h.Game().State.Chains[0]; got != 3 {
			t.Fatalf("chains = %d, want 3", got)
		}
	})
}
