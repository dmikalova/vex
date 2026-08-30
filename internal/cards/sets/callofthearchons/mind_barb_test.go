package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mind Barb
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Your opponent discards a random card from their hand.
func TestMindBarb(t *testing.T) {
	t.Run("opponent discards a random card from hand", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(MindBarb)},
			P2: ct.Side{Hand: ct.Cards(ct.Creature(), ct.Creature())},
		})

		h.P1.Play(MindBarb)

		if got := h.Game().State.Hand[1].Count; got != 1 {
			t.Fatalf("opponent hand = %d, want 1", got)
		}
		if got := h.Game().State.Discard[1].Count; got != 1 {
			t.Fatalf("opponent discard = %d, want 1", got)
		}
	})
}
