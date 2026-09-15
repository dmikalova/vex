package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Daemo-Saurus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant • Dinosaur
//
//	Play: You may exalt Daemo-Saurus -> deal 3 damage to a creature.
//	Destroyed: Steal 1 Æmber.
func TestDaemoSaurus(t *testing.T) {
	t.Run("exalting on play deals 3 damage to a creature", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, Hand: ct.Cards(DaemoSaurus)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Play(DaemoSaurus)
		h.P1.ClickCard(DaemoSaurus) // accept the may -> exalt itself
		h.P1.ClickCard(enemy)       // choose the damage target

		h.Expect(DaemoSaurus).AmberOn(1)
		h.Expect(enemy).Damage(3)
	})

	t.Run("declining the exalt deals no damage", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, Hand: ct.Cards(DaemoSaurus)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Play(DaemoSaurus)
		h.P1.ClickDone() // decline the may

		h.Expect(DaemoSaurus).AmberOn(0)
		h.Expect(enemy).Damage(0)
	})

	t.Run("steals 1 Æmber when destroyed", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, InPlay: ct.Cards(DaemoSaurus)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6)))),
				Amber:  2,
			},
		})

		h.P1.Fight(DaemoSaurus, enemy)

		h.Expect(DaemoSaurus).At(ct.Discard)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
