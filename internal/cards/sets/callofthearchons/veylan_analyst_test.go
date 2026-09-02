package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Veylan Analyst
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	After you use a card, if it is an artifact, gain 1 Æmber.
func TestVeylanAnalyst(t *testing.T) {
	t.Run("gains Æmber when you use an artifact", func(t *testing.T) {
		var analyst, relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&analyst, VeylanAnalyst),
					ct.Bind(&relic, PocketUniverse),
				),
			},
		})

		h.P1.UseAction(relic)

		h.P1.ExpectAmber(1)
	})

	t.Run("stays quiet when a creature reaps", func(t *testing.T) {
		var analyst, beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&analyst, VeylanAnalyst),
					ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(beast)

		// Only the Æmber the reap itself pays out.
		h.P1.ExpectAmber(1)
	})
}
