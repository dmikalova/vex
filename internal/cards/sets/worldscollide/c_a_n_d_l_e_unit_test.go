package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// C.A.N.D.L.E. Unit
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  1
//	Traits: Robot
//
//	After an enemy Creature reaps, draw a card.
//	Action: C.A.N.D.L.E. Unit captures 1 Æmber from your opponent.
func TestCANDLEUnit(t *testing.T) {
	t.Run("captures 1 Æmber from the opponent with its Action", func(t *testing.T) {
		var unit ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&unit, CANDLEUnit)),
			},
			P2: ct.Side{Amber: 3},
		})
		unit.Ready()

		h.P1.UseAction(CANDLEUnit)

		h.Expect(unit).AmberOn(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("draws a card after an enemy creature reaps", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(CANDLEUnit),
				Hand: ct.Cards(
					ct.Creature(), ct.Creature(), ct.Creature(),
					ct.Creature(), ct.Creature(), ct.Creature(),
				),
				Deck: ct.Cards(ct.Creature()),
			},
			P2: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Reap(foe)

		if got := len(h.Game().Hand(0)); got != 7 {
			t.Errorf("hand = %d cards, want 7", got)
		}
	})
}
