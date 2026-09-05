package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Agent Hoo-man
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Martian • Agent
//
//	Elusive.
//	Reap: Stun a friendly non-Mars creature and an enemy non-Mars creature.
func TestAgentHooMan(t *testing.T) {
	t.Run("stuns a friendly and an enemy non-mars creature", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					AgentHooMan,
					ct.Bind(&ally, ct.Creature(ct.Power(3), ct.OfHouse(card.House.Shadows))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(3), ct.OfHouse(card.House.Untamed))),
				),
			},
		})

		h.P1.Reap(AgentHooMan)

		h.Expect(ally).Stunned(true)
		h.Expect(foe).Stunned(true)
	})
}
