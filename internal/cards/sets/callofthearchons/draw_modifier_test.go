package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

func deckOf(n int) []ct.Entry {
	items := make([]any, n)
	for i := range items {
		items[i] = ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))
	}
	return ct.Cards(items...)
}

func TestMother(t *testing.T) {
	t.Run("refills its controller's hand to one additional card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(Mother), Deck: deckOf(10)},
		})

		h.P1.EndTurn()

		if got := int(h.Game().State.Hand[0].Count); got != 7 {
			t.Errorf("hand after draw = %d, want 7", got)
		}
	})
}

func TestTheHowlingPit(t *testing.T) {
	t.Run("refills each player's hand to one additional card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(TheHowlingPit), Deck: deckOf(10)},
		})

		h.P1.EndTurn()

		if got := int(h.Game().State.Hand[0].Count); got != 7 {
			t.Errorf("hand after draw = %d, want 7", got)
		}
	})
}

func TestSuccubus(t *testing.T) {
	t.Run("refills the opponent's hand to one fewer card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Succubus)},
			P2: ct.Side{Deck: deckOf(10)},
		})

		h.P1.EndTurn() // P1's turn ends; P2's turn begins
		h.P2.EndTurn() // P2 draws, reduced by Succubus

		if got := int(h.Game().State.Hand[1].Count); got != 5 {
			t.Errorf("opponent hand after draw = %d, want 5", got)
		}
	})
}
