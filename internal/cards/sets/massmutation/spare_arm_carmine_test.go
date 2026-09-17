package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Spare Arm Carmine
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Reap: Steal 1 Æmber. If you control more Mutant creatures than your opponent, steal 1 Æmber.
func TestSpareArmCarmine(t *testing.T) {
	t.Run("steals 2 with more friendly Mutants than enemy Mutants", func(t *testing.T) {
		var carmine ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					ct.Bind(&carmine, SpareArmCarmine),
					ct.Creature(ct.Traits(card.Traits.Mutant)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Reap(carmine)

		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(3)
	})

	t.Run("steals 1 when Mutants are not outnumbered", func(t *testing.T) {
		var carmine ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&carmine, SpareArmCarmine)),
			},
			P2: ct.Side{
				Amber:  5,
				InPlay: ct.Cards(ct.Creature(ct.Traits(card.Traits.Mutant))),
			},
		})

		h.P1.Reap(carmine)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(4)
	})
}
