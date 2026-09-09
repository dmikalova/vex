package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Orb of Wonder
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Special
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Orb of Wonder -> search your deck for a card and put it into your hand, then shuffle your deck.
func TestOrbOfWonder(t *testing.T) {
	var wanted, other ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Brobnar,
			InPlay: ct.Cards(OrbOfWonder),
			Deck: ct.Cards(
				ct.Bind(&wanted, ct.Creature(ct.Power(3))),
				ct.Bind(&other, ct.Creature(ct.Power(3))),
			),
		},
	})

	h.P1.UseAction(OrbOfWonder)
	h.P1.ClickCard(wanted)

	h.Expect(wanted).At(ct.Hand)
	h.Expect(other).At(ct.Deck)
	h.Expect(OrbOfWonder).At(ct.Discard)
}
