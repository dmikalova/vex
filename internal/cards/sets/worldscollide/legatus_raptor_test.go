package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Legatus Raptor
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Fight: You may exalt Legatus Raptor, and ready and use another friendly creature.
func TestLegatusRaptor(t *testing.T) {
	t.Run("exalting readies and uses another friendly creature", func(t *testing.T) {
		var raptor, ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&raptor, LegatusRaptor),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1))))},
		})
		ally.Exhaust()

		h.P1.Fight(raptor, foe)
		h.P1.ClickOption("Yes")

		h.Expect(raptor).AmberOn(1)
		h.P1.ExpectAmber(1)
	})

	t.Run("declining exalt leaves the ally exhausted", func(t *testing.T) {
		var raptor, ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&raptor, LegatusRaptor),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1))))},
		})
		ally.Exhaust()

		h.P1.Fight(raptor, foe)
		h.P1.ClickOption("No")

		h.Expect(raptor).AmberOn(0)
		h.Expect(ally).Exhausted()
	})
}
