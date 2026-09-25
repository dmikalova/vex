package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Citizen Shrix
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant
//
//	Play/Reap: Exalt Citizen Shrix. Steal 1 Æmber.
func TestCitizenShrix(t *testing.T) {
	t.Run("exalts itself and steals 1 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(CitizenShrix),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Play(CitizenShrix)

		h.Expect(CitizenShrix).AmberOn(1)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("exalts itself and steals 1 Æmber when it reaps", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(CitizenShrix),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Reap(CitizenShrix)

		h.Expect(CitizenShrix).AmberOn(1)
		h.P1.ExpectAmber(2) // 1 from reaping + 1 stolen
		h.P2.ExpectAmber(1)
	})
}
