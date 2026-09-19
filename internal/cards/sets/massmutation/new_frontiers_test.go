package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// New Frontiers
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a house. Reveal the top 3 cards of your deck. Archive each card of the chosen house and discard the others.
func TestNewFrontiers(t *testing.T) {
	t.Run("archives the chosen house and discards the others", func(t *testing.T) {
		var logos1, other, logos2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(NewFrontiers),
				Deck: ct.Cards(
					ct.Bind(&logos1, ct.Creature(ct.OfHouse(card.House.Logos))),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&logos2, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
			P2: ct.Side{},
		})

		h.P1.Play(NewFrontiers)
		h.P1.ClickOption("Logos")

		h.Expect(logos1).At(ct.Archives)
		h.Expect(logos2).At(ct.Archives)
		h.Expect(other).At(ct.Discard)
	})
}
