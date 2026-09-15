package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Killzord Mk. 9001
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains +2 power, +2 armor, and skirmish.
func TestKillzordMk9001(t *testing.T) {
	t.Run(
		"boosts its host's power and armor",
		func(t *testing.T) {
			var host ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Mars,
					InPlay: ct.Cards(
						ct.Upgraded(
							ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
							KillzordMk9001,
						),
					),
				},
			})

			h.Expect(host).Power(6)
			h.Expect(host).Armor(2)
		},
	)
}
