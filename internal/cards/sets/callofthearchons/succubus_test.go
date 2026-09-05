package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Succubus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Demon
//
//	During their "draw cards" step, your opponent refills their hand to 1 less card.
func TestSuccubus(t *testing.T) {
	t.Run("refills the opponent's hand to one fewer card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Succubus)},
			P2: ct.Side{Deck: ct.DeckOf(card.House.Logos, 10)},
		})

		h.P1.EndTurn() // P1's turn ends; P2's turn begins
		h.P2.EndTurn() // P2 draws, reduced by Succubus

		if got := int(h.Game().State.Hand[1].Count); got != 5 {
			t.Errorf("opponent hand after draw = %d, want 5", got)
		}
	})
}
