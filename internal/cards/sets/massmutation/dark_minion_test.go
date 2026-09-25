package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Dark Minion
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Mutant
//
//	Destroyed: Deal 1 damage to each enemy creature.
//	Enhance Damage.
func TestDarkMinion(t *testing.T) {
	var foe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Dis,
			InPlay: ct.Cards(DarkMinion),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
		},
	})

	h.P1.Fight(DarkMinion, foe) // the 1-power minion dies and its Destroyed fires

	h.Expect(DarkMinion).At(ct.Discard)
	h.Expect(foe).Damage(2) // 1 from the fight, 1 from Destroyed
}
