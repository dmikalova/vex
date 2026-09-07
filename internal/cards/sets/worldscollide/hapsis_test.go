package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hapsis
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Mutant • Scientist
//
//	After a creature is destroyed fighting Hapsis, ward Hapsis, and draw a card.
func TestHapsis(t *testing.T) {
	t.Run("wards itself and draws when an enemy dies fighting it", func(t *testing.T) {
		var hapsis, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&hapsis, Hapsis)),
				Deck:   ct.Cards(ct.Creature()),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2))))},
		})
		hapsis.Ready()

		h.P1.Fight(hapsis, foe)

		h.Expect(foe).At(ct.Discard)
		if !h.Game().Warded(hapsis.ID()) {
			t.Errorf("%s should be warded", hapsis.Name())
		}
		if got := len(h.Game().Hand(0)); got != 1 {
			t.Errorf("hand = %d cards, want 1", got)
		}
	})
}
