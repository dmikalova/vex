package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Even Ivan
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant • Scientist
//
//	Action: If your opponent has an even amount of Æmber, steal 1 Æmber.
func TestEvenIvan(t *testing.T) {
	t.Run("steals when the opponent has an even amount of Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(EvenIvan),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.UseAction(EvenIvan)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("does nothing when the opponent has an odd amount of Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(EvenIvan),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(EvenIvan)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
	})
}
