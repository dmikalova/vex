package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Director of Z.Y.X.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Cyborg • Scientist
//
//	Elusive.
//	At the start of your turn, archive the top card of your deck.
func TestDirectorOfZYX(t *testing.T) {
	t.Run("archives the top card of the deck at the start of the turn", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(DirectorOfZYX),
				Hand: ct.Cards(
					ct.Creature(), ct.Creature(), ct.Creature(),
					ct.Creature(), ct.Creature(), ct.Creature(),
				),
				Deck: ct.Cards(ct.Bind(&top, ct.Creature())),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)
		h.P1.ClickOption("No")

		h.Expect(top).At(ct.Archives)
	})
}
