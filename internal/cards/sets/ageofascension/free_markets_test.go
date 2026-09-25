package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Free Markets
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each house represented among cards in play, except for Sanctum, gain 1 Æmber.
func TestFreeMarkets(t *testing.T) {
	t.Run("gains 1 Æmber per house in play except Sanctum", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(FreeMarkets),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Sanctum)),
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Logos)),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Brobnar))),
			},
		})

		h.P1.Play(FreeMarkets)
		// Mars, Logos, Brobnar are represented; Sanctum is excepted.
		h.P1.ExpectAmber(3)
	})
}
