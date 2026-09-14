package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wormhole Technician
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Cyborg • Scientist
//
//	Reap: Reveal the top card of your deck. If it is a Logos card, play it. Otherwise, archive it.
func TestWormholeTechnician(t *testing.T) {
	t.Run("plays the revealed card when it is a Logos card", func(t *testing.T) {
		var tech, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&tech, WormholeTechnician)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
		})
		tech.Ready()

		h.P1.Reap(tech)

		h.Expect(top).At(ct.PlayArea)
	})

	t.Run("archives the revealed card when it is off-house", func(t *testing.T) {
		var tech, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&tech, WormholeTechnician)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})
		tech.Ready()

		h.P1.Reap(tech)

		h.Expect(top).At(ct.Archives)
	})
}
