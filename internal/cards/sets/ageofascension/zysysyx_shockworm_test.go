package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Zysysyx Shockworm
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Armor:  1
//	Traits: Martian • Soldier
//
//	After an enemy creature reaps, stun it.
func TestZysysyxShockworm(t *testing.T) {
	t.Run("stuns an enemy creature after it reaps", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ZysysyxShockworm)},
		})

		h.P1.Reap(reaper)

		h.Expect(reaper).Stunned(true)
	})

	t.Run("does not stun a friendly reaper", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ZysysyxShockworm,
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(reaper)

		h.Expect(reaper).Stunned(false)
	})
}
