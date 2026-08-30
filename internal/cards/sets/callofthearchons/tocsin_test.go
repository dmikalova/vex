package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tocsin
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Demon
//
//	Reap: Your opponent discards a random card from their hand.
func TestTocsin(t *testing.T) {
	t.Run("opponent discards a random card when it reaps", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Tocsin)},
			P2: ct.Side{Hand: ct.Cards(ct.Creature(), ct.Creature())},
		})

		h.P1.Reap(Tocsin)

		if got := h.Game().State.Hand[1].Count; got != 1 {
			t.Fatalf("opponent hand = %d, want 1", got)
		}
	})
}
