package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sci. Officer Qincan
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Alien • Proximan • Scientist
//
//	Elusive.
//	After a player chooses an active house which matches no cards in play, steal 1 Æmber.
func TestSciOfficerQincan(t *testing.T) {
	t.Run("steals 1 when the chosen house matches no card in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.StarAlliance, InPlay: ct.Cards(SciOfficerQincan)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("does nothing when the chosen house matches a card in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.StarAlliance, InPlay: ct.Cards(SciOfficerQincan)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.StarAlliance)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
	})
}
