package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Eunoia
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Beast • Cat
//
//	After a creature is destroyed in a fight with Eunoia, gain 1 Æmber. Heal 2 damage from Eunoia.
func TestEunoia(t *testing.T) {
	t.Run("gains 1 Æmber and heals 2 when an enemy dies fighting it", func(t *testing.T) {
		var eunoia, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&eunoia, Eunoia)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1))))},
		})

		h.P1.Fight(eunoia, foe)

		h.Expect(foe).At(ct.Discard)
		h.Expect(eunoia).Damage(0) // took 1 from the fight, healed 2
		h.P1.ExpectAmber(1)
	})
}
