package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tentacus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	Your opponent must give you 1 Æmber in order to use an artifact.
func TestTentacus(t *testing.T) {
	t.Run("opponent pays the controller 1 Æmber to use an artifact", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Tentacus)},
			P2: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(LibraryOfBabble),
				Deck:   ct.Cards(ct.Creature(ct.OfHouse(card.House.Logos))),
				Amber:  2,
			},
		})

		h.P1.EndTurn() // pass to the opponent
		h.P2.ChooseHouse(card.House.Logos)
		h.P2.UseAction(LibraryOfBabble)

		h.P2.ExpectAmber(1) // paid 1 of its 2
		h.P1.ExpectAmber(1) // received the toll
	})
}
