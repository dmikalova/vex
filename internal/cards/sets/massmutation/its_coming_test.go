package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// It's Coming...
//
//	House:  None
//	Type:   Tactic
//	Rarity: Special
//	Bonus:  Æmber
//
//	Play: Search your deck and discard pile for either half of a gigantic creature, reveal it, and put it into your hand. Shuffle your deck.
func TestItsComing(t *testing.T) {
	var coming, kong, plain ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			Hand:  ct.Cards(ct.Bind(&coming, engine.Rehouse(ItsComing, card.House.Logos))),
			Deck: ct.Cards(
				ct.Bind(&kong, NiffleKong),
				card.GiganticArt(NiffleKong),
				ct.Bind(&plain, ct.Creature()),
			),
		},
	})

	h.P1.Play(coming)
	h.P1.ClickCard(kong) // choose the gigantic creature to put into hand

	h.Expect(kong).At(ct.Hand)  // the chosen gigantic goes to hand
	h.Expect(plain).At(ct.Deck) // the ordinary card stays behind
}
