package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Archimedes
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg • Beast
//
//	Elusive.
//	Each neighboring creature gains, "Destroyed: Archive this creature from play."
func TestArchimedes(t *testing.T) {
	t.Run("a destroyed neighbor is archived instead of discarded", func(t *testing.T) {
		var neighbor, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					Archimedes,
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Fight(neighbor, enemy)

		h.Expect(neighbor).At(ct.Archives)
	})
}
