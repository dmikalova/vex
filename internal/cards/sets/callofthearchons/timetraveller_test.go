package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Timetraveller
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Scientist
//
//	Play: Draw 2 cards.
//	Action: Shuffle Timetraveller into its owner's deck.
func TestTimetraveller(t *testing.T) {
	t.Run("draws 2 cards when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(Timetraveller),
				Deck:  ct.Cards(ct.Creature(), ct.Creature(), ct.Creature()),
			},
		})

		h.P1.Play(Timetraveller)

		if got := len(h.Game().Hand(0)); got != 2 {
			t.Errorf("hand size = %d, want 2 (drew 2 cards)", got)
		}
	})

	t.Run("shuffles itself into the deck as an action", func(t *testing.T) {
		var tt ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&tt, Timetraveller)),
			},
		})

		h.P1.UseAction(tt)

		h.Expect(tt).At(ct.Deck)
	})
}
