package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Oubliette
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Purge a creature with power 3 or lower.
func TestOubliette(t *testing.T) {
	t.Run("purges a creature with power 3 or lower", func(t *testing.T) {
		var weak, strong ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(Oubliette)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
					ct.Bind(&strong, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Play(Oubliette)

		h.Expect(weak).At(ct.Purge)
		h.Expect(strong).At(ct.PlayArea)
	})
}
