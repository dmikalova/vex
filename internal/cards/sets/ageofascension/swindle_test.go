package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Swindle
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Alpha, Omega.
//	Play: Steal 3 Æmber.
func TestSwindle(t *testing.T) {
	t.Run("steals 3 aember when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Swindle),
				// Omega ends the turn, whose draw step recycles the discard when
				// the deck is empty; stock the deck to keep the turn deterministic.
				Deck: ct.Cards(
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(Swindle)

		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(2)
	})
}
