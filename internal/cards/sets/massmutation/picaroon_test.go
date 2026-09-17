package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Picaroon
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  0
//	Traits: Mutant • Changeling
//
//	Deploy.
//	X is the combined power of Picaroon's non-Changeling neighbors.
func TestPicaroon(t *testing.T) {
	t.Run("power is the combined power of its non-Changeling neighbors", func(t *testing.T) {
		var picaroon, changeling ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Creature(ct.Power(3)),
					ct.Bind(&picaroon, Picaroon),
					// A Changeling neighbor is excluded from the tally.
					ct.Bind(
						&changeling,
						ct.Creature(ct.Power(4), ct.Traits(card.Traits.Changeling)),
					),
				),
			},
		})

		// 3 (left, non-Changeling) + 0 (right, a Changeling is excluded) = 3.
		if got := h.Game().Power(picaroon.ID()); got != 3 {
			t.Errorf("power = %d, want 3", got)
		}
		_ = changeling
	})

	t.Run("played to an empty board while The Pale Star masks it has 1 power", func(t *testing.T) {
		var picaroon ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ThePaleStar),
				Hand:   ct.Cards(ct.Bind(&picaroon, Picaroon)),
			},
		})

		// The Pale Star's remainder-of-turn mask (1 power) is already active when
		// Picaroon deploys to the empty board: with no neighbors its own X is 0, but
		// the mask sets it to 1, so it survives rather than being destroyed at 0.
		h.P1.UseAction(ThePaleStar)
		h.P1.Play(picaroon)

		if got := h.Game().Power(picaroon.ID()); got != 1 {
			t.Errorf("masked power = %d, want 1", got)
		}
		h.Expect(picaroon).At(ct.PlayArea)
	})
}
