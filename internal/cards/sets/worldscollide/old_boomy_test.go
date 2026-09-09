package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Old Boomy
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Goblin • Scientist
//
//	Reap: Reveal cards from the top of your deck until you reveal a Brobnar card or choose to stop, archiving each card revealed this way -> deal 2 damage to Old Boomy.
func TestOldBoomy(t *testing.T) {
	t.Run("reveals until a Brobnar card, archives them, and takes 2 damage", func(t *testing.T) {
		var filler, brobnar ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(OldBoomy),
				Deck: ct.Cards(
					ct.Bind(&filler, ct.Creature(ct.OfHouse(card.House.Logos))),
					ct.Bind(&brobnar, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.Reap(OldBoomy)
		h.P1.ClickOption("Reveal another card")

		// Both revealed cards are archived; the Brobnar card ends the dig and Old
		// Boomy (power 2) takes 2 damage and is destroyed.
		h.Expect(filler).At(ct.Archives)
		h.Expect(brobnar).At(ct.Archives)
		h.Expect(OldBoomy).At(ct.Discard)
	})

	t.Run(
		"choosing to stop archives only what was revealed and deals no damage",
		func(t *testing.T) {
			var first, rest ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Brobnar,
					InPlay: ct.Cards(OldBoomy),
					Deck: ct.Cards(
						ct.Bind(&first, ct.Creature(ct.OfHouse(card.House.Logos))),
						ct.Bind(&rest, ct.Creature(ct.OfHouse(card.House.Logos))),
					),
				},
			})

			h.P1.Reap(OldBoomy)
			h.P1.ClickOption("Stop")

			h.Expect(first).At(ct.Archives)
			h.Expect(rest).At(ct.Deck)
			h.Expect(OldBoomy).At(ct.PlayArea)
		},
	)
}
