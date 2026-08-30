package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Subtle Maul
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Action: Your opponent discards a random card from their hand.
func TestSubtleMaul(t *testing.T) {
	t.Run("opponent discards a random card as an action", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(SubtleMaul)},
			P2: ct.Side{Hand: ct.Cards(ct.Creature(), ct.Creature())},
		})

		h.P1.UseAction(SubtleMaul)

		if got := h.Game().State.Hand[1].Count; got != 1 {
			t.Fatalf("opponent hand = %d, want 1", got)
		}
	})
}
