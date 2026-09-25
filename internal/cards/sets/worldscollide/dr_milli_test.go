package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Dr. Milli
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Scientist
//
//	Play: For each creature your opponent controls in excess of you, not counting Dr. Milli, archive a card from your hand.
func TestDrMilli(t *testing.T) {
	t.Run(
		"archives one card per opponent creature in excess, not counting itself",
		func(t *testing.T) {
			var a, b, c ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Logos,
					Hand: ct.Cards(
						DrMilli,
						ct.Bind(&a, ct.Creature()),
						ct.Bind(&b, ct.Creature()),
						ct.Bind(&c, ct.Creature()),
					),
				},
				P2: ct.Side{
					House:  card.House.Brobnar,
					InPlay: ct.Cards(ct.Creature(), ct.Creature(), ct.Creature()),
				},
			})

			// P1 has 0 creatures before Dr. Milli; the opponent has 3. Not counting
			// Dr. Milli, the excess is 3, so three cards are archived.
			h.P1.Play(DrMilli)
			h.P1.ClickCard(a)
			h.P1.ClickCard(b)
			// The third archive has only one candidate left and resolves automatically.

			h.Expect(a).At(ct.Archives)
			h.Expect(b).At(ct.Archives)
			h.Expect(c).At(ct.Archives)
		},
	)

	t.Run("archives nothing when the opponent has no excess", func(t *testing.T) {
		var a ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Creature(), ct.Creature()),
				Hand:   ct.Cards(DrMilli, ct.Bind(&a, ct.Creature())),
			},
			P2: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Creature(), ct.Creature()),
			},
		})

		// P1 has 2 creatures (3 with Dr. Milli, 2 not counting it); opponent has 2.
		// Excess is 0, so nothing is archived.
		h.P1.Play(DrMilli)

		h.Expect(a).At(ct.Hand)
	})
}
