package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Tomes Gigantica
//
//	House:  None
//	Type:   Tactic
//	Rarity: Special
//	Bonus:  Æmber
//
//	Play: Search your deck and discard pile for two halves of a gigantic creature, reveal them, and put them into your hand. Shuffle your deck. Purge Tomes Gigantica.
func TestTomesGigantica(t *testing.T) {
	var tomes, base, art, plain ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			Hand:  ct.Cards(ct.Bind(&tomes, engine.Rehouse(TomesGigantica, card.House.Logos))),
			Deck: ct.Cards(
				ct.Bind(&base, NiffleKong),
				ct.Bind(&art, card.GiganticArt(NiffleKong)),
				ct.Bind(&plain, ct.Creature()),
			),
		},
	})

	h.P1.Play(tomes)
	h.P1.ClickCard(base) // take the first half
	h.P1.ClickCard(art)  // take the second half

	h.Expect(base).At(ct.Hand) // both halves go to hand
	h.Expect(art).At(ct.Hand)
	h.Expect(plain).At(ct.Deck)  // the ordinary card stays behind
	h.Expect(tomes).At(ct.Purge) // Tomes purges itself instead of discarding
}
