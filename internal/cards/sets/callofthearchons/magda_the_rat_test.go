package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Magda the Rat
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Elf • Thief
//
//	Elusive.
//	Play: Steal 2 Æmber.
//	Leaves Play: Your opponent steals 2 Æmber.
func TestMagdaTheRat(t *testing.T) {
	t.Run("steals on the way in and gives it back on the way out", func(t *testing.T) {
		var magda, killer, decoy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ct.Bind(&magda, MagdaTheRat)),
			},
			P2: ct.Side{
				Amber: 5,
				InPlay: ct.Cards(
					ct.Bind(&killer, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(9))),
					ct.Bind(&decoy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(1))),
				),
			},
		})

		h.P1.Play(magda)
		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(3)

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		// Elusive soaks the first attack each turn, so the decoy spends it.
		h.P2.Fight(decoy, magda)
		h.P2.Fight(killer, magda)

		h.Expect(magda).At(ct.Discard)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(5)
	})
}
