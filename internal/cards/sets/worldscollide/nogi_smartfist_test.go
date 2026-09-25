package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Nogi Smartfist
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Giant • Scientist
//
//	Fight: Draw 2 cards. Discard 2 random cards from your hand.
func TestNogiSmartfist(t *testing.T) {
	t.Run("draws 2 then discards 2 random cards when it fights", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(NogiSmartfist),
				Hand:   ct.Cards(ct.Tactic(), ct.Tactic()),
				Deck:   ct.Cards(ct.Tactic(), ct.Tactic()),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.Fight(NogiSmartfist, foe)

		// 2 in hand + 2 drawn - 2 discarded = 2 in hand; the 2 discarded pile up.
		if got := h.Game().State.Hand[0].Count; got != 2 {
			t.Errorf("hand count = %d, want 2", got)
		}
		if got := h.Game().State.Discard[0].Count; got != 2 {
			t.Errorf("discard count = %d, want 2", got)
		}
		if got := h.Game().State.Deck[0].Count; got != 0 {
			t.Errorf("deck count = %d, want 0", got)
		}
	})
}
