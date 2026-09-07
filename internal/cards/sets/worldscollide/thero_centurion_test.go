package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Thero Centurion
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Play/Fight: Thero Centurion captures 1 Æmber from your opponent.
func TestTheroCenturion(t *testing.T) {
	t.Run("captures 1 Æmber from the opponent when played", func(t *testing.T) {
		var thero ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, Hand: ct.Cards(ct.Bind(&thero, TheroCenturion))},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(TheroCenturion)

		h.Expect(thero).AmberOn(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("captures 1 more Æmber when it fights", func(t *testing.T) {
		var thero, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&thero, TheroCenturion)),
			},
			P2: ct.Side{
				Amber:  3,
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2)))),
			},
		})
		thero.Ready()

		h.P1.Fight(thero, foe)

		h.Expect(thero).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
