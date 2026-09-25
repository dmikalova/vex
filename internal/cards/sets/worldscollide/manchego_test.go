package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Manchego
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Thief
//
//	Play: If you have 5 or fewer cards in your deck, steal 2 Æmber.
//	Fight/Reap: You may shuffle Manchego into its owner's deck.
func TestManchego(t *testing.T) {
	t.Run("steals with a small deck", func(t *testing.T) {
		var m ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ct.Bind(&m, Manchego)),
				Deck:  ct.Cards(ct.Creature(), ct.Creature()),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(m)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})

	t.Run("does not steal with a large deck", func(t *testing.T) {
		var m ct.Card
		deck := make([]ct.Entry, 0, 6)
		for range 6 {
			deck = append(deck, ct.Bind(nil, ct.Creature()))
		}
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ct.Bind(&m, Manchego)),
				Deck:  deck,
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(m)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
	})

	t.Run("reaps to shuffle itself into the deck", func(t *testing.T) {
		var m ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&m, Manchego)),
			},
		})
		m.Ready()

		h.P1.Reap(m)
		h.P1.ClickCard(m)

		h.Expect(m).At(ct.Deck)
	})
}
