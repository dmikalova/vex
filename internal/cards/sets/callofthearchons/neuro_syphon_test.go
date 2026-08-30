package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Neuro Syphon
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If your opponent has more Æmber than you, steal 1 Æmber, and draw a card.
func TestNeuroSyphon(t *testing.T) {
	t.Run("steals 1 Æmber and draws when the opponent is ahead", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(NeuroSyphon),
				Deck:  ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(NeuroSyphon)

		h.P2.ExpectAmber(4) // 1 stolen
		h.Expect(top).At(ct.Hand)
	})

	t.Run("does nothing when not behind", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(NeuroSyphon)},
			P2: ct.Side{Amber: 0},
		})

		h.P1.Play(NeuroSyphon)

		h.P2.ExpectAmber(0)
	})
}
