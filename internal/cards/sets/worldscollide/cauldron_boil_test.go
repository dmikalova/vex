package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cauldron Boil
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to each creature for each point of damage on it.
func TestCauldronBoil(t *testing.T) {
	var big, unhurt, doomed ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Untamed,
			Hand:   ct.Cards(CauldronBoil),
			InPlay: ct.Cards(ct.Bind(&big, ct.Creature(ct.Power(10)))),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&unhurt, ct.Creature(ct.Power(5))),
				ct.Bind(&doomed, ct.Creature(ct.Power(4))),
			),
		},
	})

	big.Damaged(3)
	doomed.Damaged(3)

	h.P1.Play(CauldronBoil)

	h.Expect(big).Damage(6)         // 3 already on it, +3 for the 3 points there
	h.Expect(unhurt).Damage(0)      // undamaged, so it takes nothing
	h.Expect(doomed).At(ct.Discard) // 3 + 3 = 6 >= 4 power, destroyed
	h.P1.ExpectAmber(1)             // from Cauldron Boil's Æmber bonus
}
