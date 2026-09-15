package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Martyr's End
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy any number of friendly creatures. For each creature destroyed this way, gain 1 Æmber.
func TestMartyrsEnd(t *testing.T) {
	var a, b ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Sanctum,
			InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Sanctum))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Sanctum))),
			),
			Hand: ct.Cards(MartyrsEnd),
		},
	})

	h.P1.Play(MartyrsEnd)
	h.P1.ClickCard(a)
	h.P1.ClickCard(b)

	h.Expect(a).At(ct.Discard)
	h.Expect(b).At(ct.Discard)
	// One Æmber pip from playing Martyr's End plus one for each creature destroyed.
	h.P1.ExpectAmber(3)
}
