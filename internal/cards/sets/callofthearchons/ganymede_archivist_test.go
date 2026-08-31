package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ganymede Archivist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Scientist
//
//	Reap: Archive a card from your hand.
func TestGanymedeArchivist(t *testing.T) {
	t.Run("archives a card from hand when it reaps", func(t *testing.T) {
		var spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					ct.Bind(&spare, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
				),
				InPlay: ct.Cards(GanymedeArchivist),
			},
		})

		h.P1.Reap(GanymedeArchivist)

		h.Expect(spare).At(ct.Archives)
	})
}
