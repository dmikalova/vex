package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dendrix
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Fight: Your opponent discards a random card from their hand.
func TestDendrix(t *testing.T) {
	t.Run("opponent discards a random card when Dendrix fights", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Dendrix),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
				Hand:   ct.Cards(ct.Creature(), ct.Creature()),
			},
		})

		h.P1.Fight(Dendrix, foe)

		if got := len(h.Game().Hand(1)); got != 1 {
			t.Errorf("opponent hand = %d, want 1", got)
		}
	})
}
