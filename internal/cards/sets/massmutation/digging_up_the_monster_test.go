package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Digging Up the Monster
//
//	House:  None
//	Type:   Tactic
//	Rarity: Special
//	Bonus:  Æmber
//
//	Play: Search your deck and discard pile for two halves of a gigantic creature, reveal them, shuffle your deck, and put them into the top of your deck.
func TestDiggingUpTheMonster(t *testing.T) {
	var dig, base, art, plain ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			Hand:  ct.Cards(ct.Bind(&dig, engine.Rehouse(DiggingUpTheMonster, card.House.Logos))),
			Discard: ct.Cards(
				ct.Bind(&base, NiffleKong),
				ct.Bind(&art, card.GiganticArt(NiffleKong)),
				ct.Bind(&plain, ct.Creature()),
			),
		},
	})

	h.P1.Play(dig)
	h.P1.ClickCard(base) // take the first half
	h.P1.ClickCard(art)  // take the second half

	h.Expect(base).At(ct.Deck) // both halves move onto the deck
	h.Expect(art).At(ct.Deck)
	h.Expect(plain).At(ct.Discard) // the ordinary card stays in the discard pile
}
