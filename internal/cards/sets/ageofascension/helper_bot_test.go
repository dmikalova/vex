package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Helper Bot
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Robot
//
//	Play: Play a non-Logos card.
func TestHelperBot(t *testing.T) {
	t.Run("plays a non-Logos card from hand", func(t *testing.T) {
		var brobnar ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					HelperBot,
					ct.Bind(&brobnar, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.Play(HelperBot)
		h.Expect(brobnar).At(ct.PlayArea)
	})

	t.Run("a Logos card cannot be chosen", func(t *testing.T) {
		var logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					HelperBot,
					ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Play(HelperBot)
		h.Expect(logos).At(ct.Hand)
	})
}
