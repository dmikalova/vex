package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Storm Crawler
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Robot
//
//	Storm Crawler deals 1 Damage when fighting.
//	After an enemy creature reaps, stun it.
func TestStormCrawler(t *testing.T) {
	t.Run("deals only 1 fight damage", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(StormCrawler),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(StormCrawler, foe)

		h.Expect(foe).Damage(1)
	})

	t.Run("stuns an enemy creature after it reaps", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(StormCrawler)},
		})

		h.P1.Reap(reaper)

		h.Expect(reaper).Stunned(true)
	})
}
