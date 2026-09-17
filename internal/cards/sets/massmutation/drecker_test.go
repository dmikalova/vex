package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Drecker
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Imp
//
//	Damage dealt to Drecker's neighbors during fights is also dealt to Drecker.
//	Reap: Steal 1 Æmber.
func TestDrecker(t *testing.T) {
	t.Run("shares a neighbor's fight damage", func(t *testing.T) {
		var neighbor, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					Drecker,
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(6))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		// The neighbor takes 3 retaliation from the foe it fights; that same damage
		// is also dealt to Drecker.
		h.P1.Fight(neighbor, foe)
		h.Expect(neighbor).Damage(3)
		h.Expect(Drecker).Damage(3)
	})

	t.Run("Reap: Steal 1", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Drecker),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Reap(Drecker)
		h.P1.ExpectAmber(2) // 1 reap + 1 stolen
		h.P2.ExpectAmber(1)
	})
}
