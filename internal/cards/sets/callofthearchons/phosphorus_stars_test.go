package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Phosphorus Stars
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Stun each non-Mars creature. Gain 2 chains.
func TestPhosphorusStars(t *testing.T) {
	t.Run("stuns each non-Mars creature and gains 2 chains", func(t *testing.T) {
		var marsAlly, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				Hand:   ct.Cards(PhosphorusStars),
				InPlay: ct.Cards(ct.Bind(&marsAlly, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
			},
		})

		h.P1.Play(PhosphorusStars)

		h.Expect(foe).Stunned(true)
		h.Expect(marsAlly).Stunned(false)
		if got := h.Game().State.Chains[0]; got != 2 {
			t.Fatalf("chains = %d, want 2", got)
		}
	})
}
