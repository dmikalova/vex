package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Build Your Champion
//
//	House:  None
//	Type:   Tactic
//	Rarity: Special
//	Bonus:  Æmber
//
//	Play: Search your deck and discard pile for two halves of a gigantic creature, reveal them, and put them into your archives. Shuffle your deck.
func TestBuildYourChampion(t *testing.T) {
	var champ, base, art, plain ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			Hand:  ct.Cards(ct.Bind(&champ, engine.Rehouse(BuildYourChampion, card.House.Logos))),
			Deck: ct.Cards(
				ct.Bind(&base, NiffleKong),
				ct.Bind(&art, card.GiganticArt(NiffleKong)),
				ct.Bind(&plain, ct.Creature()),
			),
		},
	})

	h.P1.Play(champ)
	h.P1.ClickCard(base) // take the first half
	h.P1.ClickCard(art)  // take the second half

	h.Expect(base).At(ct.Archives) // both halves are archived
	h.Expect(art).At(ct.Archives)
	h.Expect(plain).At(ct.Deck) // the ordinary card stays behind
}
