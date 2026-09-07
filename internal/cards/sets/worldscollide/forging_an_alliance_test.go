package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Forging an Alliance
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Forge a key at +7 Æmber current cost, reduced by 1 Æmber for each house represented among cards in play (to a maximum of 6).
func TestForgingAnAlliance(t *testing.T) {
	t.Run("forges at +7 reduced by 1 per house among cards in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ForgingAnAlliance),
				Amber: 12,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Logos)),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Brobnar))),
			},
		})

		// Mars, Logos, Brobnar = 3 houses, so +7 reduced by 3 = +4, plus the
		// base 6 key cost = 10, plus the card's own 1 Æmber bonus.
		h.P1.Play(ForgingAnAlliance)
		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(3)
	})
}
