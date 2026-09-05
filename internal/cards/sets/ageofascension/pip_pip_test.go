package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pip Pip
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Scientist
//
//	After an enemy creature reaps, stun it.
func TestPipPip(t *testing.T) {
	t.Run("stuns an enemy creature after it reaps", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(PipPip)},
		})

		h.P1.Reap(reaper)

		h.Expect(reaper).Stunned(true)
	})
}
