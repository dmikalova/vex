package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Psychic Bug
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Æmber:  1
//	Traits: Cyborg • Insect
//
//	Play/Reap: Reveal your opponent's hand.
func TestPsychicBug(t *testing.T) {
	t.Run("reveals the opponent's hand when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(PsychicBug)},
			P2: ct.Side{Hand: ct.Cards(ct.Creature(), ct.Creature())},
		})

		h.P1.Play(PsychicBug)

		// Revealing is informational; the hand is untouched.
		h.Expect(PsychicBug).At(ct.PlayArea)
		if got := h.Game().State.Hand[1].Count; got != 2 {
			t.Fatalf("opponent hand = %d, want 2 (reveal does not discard)", got)
		}
	})
}
