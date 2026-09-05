package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Z.Y.X. Researcher
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Scientist
//
//	Play: Choose one:
//	- Archive the top card of your deck
//	- Archive the top card of your discard pile.
func TestZYXResearcher(t *testing.T) {
	t.Run("archives the top card of the deck", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(ZYXResearcher),
				Deck:  ct.Cards(ct.Bind(&top, ct.Creature())),
			},
			P2: ct.Side{},
		})

		h.P1.Play(ZYXResearcher)
		h.P1.ClickOption("deck")

		h.Expect(top).At(ct.Archives)
	})

	t.Run("archives the top card of the discard pile", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Logos,
				Hand:    ct.Cards(ZYXResearcher),
				Discard: ct.Cards(ct.Bind(&top, ct.Creature())),
			},
			P2: ct.Side{},
		})

		h.P1.Play(ZYXResearcher)
		h.P1.ClickOption("discard")

		h.Expect(top).At(ct.Archives)
	})
}
