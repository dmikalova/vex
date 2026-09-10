package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Malison
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Fight: Move an enemy creature anywhere in its controller's battleline -> if it is on a flank, the chosen creature captures 1 Æmber from your opponent.
func TestMalison(t *testing.T) {
	var defender, mid, right ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Malison)},
		P2: ct.Side{
			Amber: 5,
			InPlay: ct.Cards(
				ct.Bind(&defender, ct.Creature(ct.Power(1))),
				ct.Bind(&mid, ct.Creature(ct.Power(9))),
				ct.Bind(&right, ct.Creature(ct.Power(9))),
			),
		},
	})

	h.P1.Fight(Malison, defender)
	h.P1.ClickCard(mid)            // move the interior creature
	h.P1.ClickOption("Left flank") // onto a flank

	// The moved creature ends on a flank, so it captures 1 Æmber from its own side.
	h.Expect(mid).AmberOn(1)
	h.P2.ExpectAmber(4)
	_ = right
}
