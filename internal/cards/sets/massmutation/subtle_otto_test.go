package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Subtle Otto
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Mutant • Thief
//
//	Play: Your opponent discards a random card from their hand.
func TestSubtleOtto(t *testing.T) {
	t.Run("opponent discards a random card from their hand when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(SubtleOtto)},
			P2: ct.Side{Hand: ct.Cards(ct.Creature(), ct.Creature())},
		})

		h.P1.Play(SubtleOtto)

		if got := h.Game().State.Hand[1].Count; got != 1 {
			t.Fatalf("opponent hand = %d, want 1", got)
		}
	})
}
