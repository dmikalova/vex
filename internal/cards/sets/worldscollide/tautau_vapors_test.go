package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tautau Vapors
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Draw 2 cards. Archive a card from your hand.
func TestTautauVapors(t *testing.T) {
	t.Run("draws 2 cards then archives one", func(t *testing.T) {
		var top1, top2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(TautauVapors),
				Deck:  ct.Cards(ct.Bind(&top1, ct.Creature()), ct.Bind(&top2, ct.Creature())),
			},
		})

		h.P1.Play(TautauVapors)
		h.P1.ClickCard(top1)

		h.Expect(top1).At(ct.Archives)
		h.Expect(top2).At(ct.Hand)
	})
}
