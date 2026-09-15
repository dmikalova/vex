package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Foozle
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Reap: If an enemy creature has been destroyed this turn, gain 1 Æmber.
func TestFoozle(t *testing.T) {
	t.Run("gains the extra Æmber once an enemy creature is destroyed", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					Foozle,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Fight(ally, foe) // destroys the enemy creature this turn
		h.P1.Reap(Foozle)

		h.P1.ExpectAmber(2) // 1 for reaping, 1 for the condition
	})

	t.Run("gains only the reap Æmber when no enemy creature was destroyed", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Foozle),
			},
			P2: ct.Side{},
		})

		h.P1.Reap(Foozle)

		h.P1.ExpectAmber(1)
	})
}
