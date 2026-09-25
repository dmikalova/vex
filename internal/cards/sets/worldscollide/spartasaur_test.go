package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Spartasaur
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	After a friendly creature is destroyed, destroy each non-Dinosaur creature.
//	Fight: Gain 2 Æmber.
func TestSpartasaur(t *testing.T) {
	t.Run("a friendly death destroys each non-Dinosaur creature", func(t *testing.T) {
		var victim, dinoAlly, nonDinoAlly, foughtEnemy, nonDinoEnemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					Spartasaur,
					ct.Bind(&victim, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(2))),
					ct.Bind(&dinoAlly, ct.Creature(ct.Traits(card.Traits.Dinosaur))),
					ct.Bind(&nonDinoAlly, ct.Creature(ct.Traits(card.Traits.Human))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foughtEnemy, ct.Creature(ct.Power(5), ct.Traits(card.Traits.Dinosaur))),
				ct.Bind(&nonDinoEnemy, ct.Creature(ct.Traits(card.Traits.Human))),
			)},
		})

		h.P1.Fight(victim, foughtEnemy)

		h.Expect(victim).At(ct.Discard)
		h.Expect(nonDinoAlly).At(ct.Discard)
		h.Expect(nonDinoEnemy).At(ct.Discard)
		h.Expect(dinoAlly).At(ct.PlayArea)
		h.Expect(foughtEnemy).At(ct.PlayArea)
		h.Expect(Spartasaur).At(ct.PlayArea)
	})

	t.Run("Fight gains 2 Æmber", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(Spartasaur),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(1))))},
		})

		h.P1.Fight(Spartasaur, enemy)

		h.P1.ExpectAmber(2)
	})
}
