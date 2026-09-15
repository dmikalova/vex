package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Keyforgery
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	When your opponent would forge a key, they name a house. Reveal a random card from your hand. If that card is not of the named house, destroy Keyforgery, and they do not forge that key.
func TestKeyforgery(t *testing.T) {
	t.Run("a wrong guess destroys Keyforgery and prevents the forge", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(Keyforgery),
				Hand:   ct.Cards(ct.Creature(ct.OfHouse(card.House.Logos))),
			},
			P2: ct.Side{Amber: 6},
		})

		h.P1.EndTurn()              // P2 forges at the start of their turn; Keyforgery interrupts
		h.P2.ClickOption("Brobnar") // the revealed card is Logos, not Brobnar

		h.Expect(Keyforgery).At(ct.Discard)
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(6)
	})

	t.Run(
		"a correct guess leaves Keyforgery in play and lets the forge happen",
		func(t *testing.T) {
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Shadows,
					InPlay: ct.Cards(Keyforgery),
					Hand:   ct.Cards(ct.Creature(ct.OfHouse(card.House.Logos))),
				},
				P2: ct.Side{Amber: 6},
			})

			h.P1.EndTurn()
			h.P2.ClickOption("Logos") // matches the revealed card's house

			h.Expect(Keyforgery).At(ct.PlayArea)
			h.P2.ExpectKeys(1)
			h.P2.ExpectAmber(0)
		},
	)
}
