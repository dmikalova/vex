package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// ANT1-10NY
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Robot
//
//	At the end of your turn, move 1 Æmber from ANT1-10NY to your opponent's pool.
//	Play: ANT1-10NY captures all your opponent's Æmber.
func TestANT110NY(t *testing.T) {
	t.Run("captures all of the opponent's Æmber, then feeds 1 back each turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ANT110NY),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(ANT110NY)

		h.Expect(ANT110NY).AmberOn(3)
		h.P2.ExpectAmber(0)

		h.P1.EndTurn() // end-of-turn moves 1 back to the opponent's pool

		h.Expect(ANT110NY).AmberOn(2)
		h.P2.ExpectAmber(1)
	})

	t.Run("captures nothing when the opponent has no Æmber", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ANT110NY),
			},
		})

		h.P1.Play(ANT110NY)

		h.Expect(ANT110NY).AmberOn(0)
	})
}
