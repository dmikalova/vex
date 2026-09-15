package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bloodshard Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	After a creature reaps, destroy it.
func TestBloodshardImp(t *testing.T) {
	t.Run("destroys whatever creature just reaped", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					BloodshardImp,
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(reaper)

		h.Expect(reaper).At(ct.Discard)
	})
}
