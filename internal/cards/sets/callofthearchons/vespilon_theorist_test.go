package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Vespilon Theorist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	Elusive.
//	Reap: Choose a house - reveal the top card of your deck. If it is of the chosen house, archive the top card of your deck, and gain 1 Æmber. Otherwise, discard the top card of your deck.
func TestVespilonTheorist(t *testing.T) {
	t.Run("archives the revealed card of the chosen house", func(t *testing.T) {
		var theorist, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&theorist, VespilonTheorist)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})
		theorist.Ready()

		h.P1.Reap(theorist)
		h.P1.ClickOption("Mars")

		h.Expect(top).At(ct.Archives)
		h.P1.ExpectAmber(2) // one for the reap, one for the match
	})

	t.Run("discards the revealed card of another house", func(t *testing.T) {
		var theorist, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&theorist, VespilonTheorist)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})
		theorist.Ready()

		h.P1.Reap(theorist)
		h.P1.ClickOption("Shadows")

		h.Expect(top).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})
}
