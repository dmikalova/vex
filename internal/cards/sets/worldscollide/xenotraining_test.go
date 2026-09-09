package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Xenotraining
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each house represented among friendly creatures, a friendly creature captures 1 Æmber from your opponent.
func TestXenotraining(t *testing.T) {
	var a, b ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.StarAlliance,
			Hand:  ct.Cards(Xenotraining),
			InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Saurian))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Logos))),
			),
		},
		P2: ct.Side{Amber: 5},
	})

	// Two distinct houses among friendly creatures → two captures.
	h.P1.Play(Xenotraining)
	h.P1.ClickCard(a)
	h.P1.ClickCard(a)

	h.Expect(a).AmberOn(2)
	h.Expect(b).AmberOn(0)
	h.P2.ExpectAmber(3)
}
