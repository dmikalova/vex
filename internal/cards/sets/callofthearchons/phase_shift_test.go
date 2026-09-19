package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Phase Shift
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Play a non-Logos card.
func TestPhaseShift(t *testing.T) {
	t.Run("plays a non-Logos card from hand", func(t *testing.T) {
		var brobnar ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					PhaseShift,
					ct.Bind(&brobnar, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.Play(PhaseShift)
		h.Expect(brobnar).At(ct.PlayArea)
	})

	t.Run("a Logos card cannot be chosen", func(t *testing.T) {
		var logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					PhaseShift,
					ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Play(PhaseShift)
		h.Expect(logos).At(ct.Hand)
	})

	// The first-turn rule limits the first player to one play of their own
	// volition, but a card another card plays is not limited: Phase Shift is P1's
	// one play, yet the non-Logos card it forces still enters play.
	t.Run("forces a play on the first player's first turn", func(t *testing.T) {
		var brobnar ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					PhaseShift,
					ct.Bind(&brobnar, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})
		h.Game().State.FirstTurnPlayLimit[0] = true

		h.P1.Play(PhaseShift)
		h.Expect(brobnar).At(ct.PlayArea)
	})
}
