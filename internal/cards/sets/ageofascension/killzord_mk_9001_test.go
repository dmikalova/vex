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
//	This creature gains +2 power and +2 armor and skirmish.
//	This creature gains, "Fight: Gain 1 chain."
func TestKillzordMk9001(t *testing.T) {
	t.Run(
		"boosts its host and gives its controller a chain when the host fights",
		func(t *testing.T) {
			var host, foe ct.Card
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
				P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
			})

			h.Expect(host).Power(6)
			h.Expect(host).Armor(2)

			h.P1.Fight(host, foe)

			if got := h.Game().State.Chains[0]; got != 1 {
				t.Errorf("controller should gain 1 chain, got %d", got)
			}
		},
	)
}
