package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gabos Longarms
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Demon
//
//	Before Fight: Choose a creature - Gabos Longarms deals its fight damage to the chosen creature instead of to the creature it is fighting.
func TestGabosLongarms(t *testing.T) {
	var def, bystander ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(GabosLongarms)},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&def, ct.Creature(ct.Power(3))),
			ct.Bind(&bystander, ct.Creature(ct.Power(2))),
		)},
	})

	h.P1.Fight(GabosLongarms, def)
	h.P1.ClickCard(bystander) // Before Fight: send Gabos's damage to the bystander

	h.Expect(bystander).At(ct.Discard) // took 5, power 2 -> destroyed
	h.Expect(def).Damage(0)            // the creature it fought takes none
	h.Expect(GabosLongarms).Damage(3)  // Gabos still takes the defender's 3 back
}
