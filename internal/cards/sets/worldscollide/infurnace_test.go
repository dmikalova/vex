package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Infurnace
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Play: Purge 2 cards from a discard pile. For each Æmber bonus icon on the purged cards, your opponent loses 1 Æmber.
func TestInfurnace(t *testing.T) {
	var a, b ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(Infurnace)},
		P2: ct.Side{
			Amber: 5,
			Discard: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.AemberBonus(2))),
				ct.Bind(&b, ct.Creature(ct.AemberBonus(1))),
			),
		},
	})

	h.P1.Play(Infurnace)
	h.P1.ClickCard(a)

	h.Expect(a).At(ct.Purge)
	h.Expect(b).At(ct.Purge)
	h.P2.ExpectAmber(2) // 5 - (2 + 1) Æmber bonus icons
}
