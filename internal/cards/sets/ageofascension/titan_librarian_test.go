package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Titan Librarian
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	At the end of your turn, if Titan Librarian is not on a flank, archive a card from your hand.
func TestTitanLibrarian(t *testing.T) {
	t.Run("archives a card at end of turn while not on a flank", func(t *testing.T) {
		var buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Creature(ct.Power(2)),
					TitanLibrarian,
					ct.Creature(ct.Power(2)),
				),
				Hand: ct.Cards(ct.Bind(&buried, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.EndTurn()

		h.Expect(buried).At(ct.Archives)
	})

	t.Run("archives nothing while on a flank", func(t *testing.T) {
		var buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					TitanLibrarian,
					ct.Creature(ct.Power(2)),
				),
				Hand: ct.Cards(ct.Bind(&buried, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.EndTurn()

		h.Expect(buried).At(ct.Hand)
	})
}
