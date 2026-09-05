package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Aemberspine Mongrel
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Hazardous 3.
//	After an enemy creature reaps, gain 1 Æmber.
func TestAemberspineMongrel(t *testing.T) {
	t.Run("gains Æmber when an enemy creature reaps", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(AemberspineMongrel)},
		})

		h.P1.Reap(reaper)

		h.P2.ExpectAmber(1)
	})
}
