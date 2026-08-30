package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Blinding Light
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a house - stun each creature of the chosen house.
func TestBlindingLight(t *testing.T) {
	t.Run("stuns each creature of the chosen house, sparing the rest", func(t *testing.T) {
		var marsFoe, shadowFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(BlindingLight),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&marsFoe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&shadowFoe, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
				),
			},
		})

		h.P1.Play(BlindingLight)
		h.P1.ExpectPrompt("Choose a house").Source("Blinding Light")
		h.P1.ClickOption("Mars")

		h.Expect(marsFoe).Stunned(true)
		h.Expect(shadowFoe).Stunned(false)
	})
}
