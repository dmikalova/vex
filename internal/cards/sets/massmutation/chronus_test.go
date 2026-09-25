package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Chronus
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant
//
//	After you resolve a Draw bonus icon, you may archive a card from your hand.
//	Enhance Draw Draw.
func TestChronus(t *testing.T) {
	t.Run("may archive a card after a Draw bonus icon resolves", func(t *testing.T) {
		var bearer, toArchive ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Chronus),
				Hand: ct.Cards(
					ct.Bind(
						&bearer,
						ct.Creature(ct.OfHouse(card.House.Logos), ct.Bonus(card.Bonus.Draw)),
					),
					ct.Bind(&toArchive, ct.Tactic(ct.OfHouse(card.House.Logos))),
				),
				Deck: ct.Cards(ct.Creature()),
			},
		})

		h.P1.Play(bearer)         // the Draw bonus icon draws, then Chronus reacts
		h.P1.ClickCard(Chronus)   // click Chronus to accept its offer
		h.P1.ClickCard(toArchive) // archive this card from hand

		h.Expect(toArchive).At(ct.Archives)
	})
}
