package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Orb of Invidius
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	After a creature reaps, stun it.
func TestOrbOfInvidius(t *testing.T) {
	t.Run("stuns whatever creature just reaped", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					OrbOfInvidius,
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(reaper)

		h.Expect(reaper).Stunned(true)
	})
}
