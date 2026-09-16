package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Scout Pete
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Alien
//
//	Play/Fight/Reap: Look at the top card of your deck and you may discard that card.
func TestScoutPete(t *testing.T) {
	t.Run("discards the looked-at card when the controller chooses to", func(t *testing.T) {
		var pete, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&pete, ScoutPete)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.StarAlliance)))),
			},
		})

		h.P1.Reap(pete)
		h.P1.ClickCard(top)

		h.Expect(top).At(ct.Discard)
	})

	t.Run("keeps the card on top when the controller declines", func(t *testing.T) {
		var pete, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&pete, ScoutPete)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.StarAlliance)))),
			},
		})

		h.P1.Reap(pete)
		h.P1.ClickDone()

		h.Expect(top).At(ct.Deck)
	})
}
