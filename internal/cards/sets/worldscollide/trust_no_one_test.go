package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Trust No One
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose one:
//	- If there are no friendly Creatures in play, for each house represented among enemy Creatures, steal 1 Æmber
//	- Steal 1 Æmber.
func TestTrustNoOne(t *testing.T) {
	t.Run("steals 1 Æmber while you control creatures", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Hand:   ct.Cards(TrustNoOne),
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Shadows))),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(TrustNoOne)
		h.P1.ClickOption("steal 1 Æmber")
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(4)
	})

	t.Run(
		"with no friendly creatures, steals 1 per house among enemy creatures",
		func(t *testing.T) {
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Shadows,
					Hand:  ct.Cards(TrustNoOne),
				},
				P2: ct.Side{
					Amber: 5,
					InPlay: ct.Cards(
						ct.Creature(ct.OfHouse(card.House.Mars)),
						ct.Creature(ct.OfHouse(card.House.Logos)),
						ct.Creature(ct.OfHouse(card.House.Brobnar)),
						ct.Creature(ct.OfHouse(card.House.Untamed)),
					),
				},
			})

			// Mars, Logos, Brobnar, Untamed = 4 houses, uncapped.
			h.P1.Play(TrustNoOne)
			h.P1.ClickOption("for each house")
			h.P1.ExpectAmber(4)
			h.P2.ExpectAmber(1)
		},
	)
}
