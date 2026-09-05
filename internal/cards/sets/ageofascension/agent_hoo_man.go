package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

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
var AgentHooMan = card.New(
	"Agent Hoo-man",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 160),
	card.WithPower(2),
	card.WithTraits(card.Traits.Martian, card.Traits.Agent),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.Sequence{
			Effects: []card.Effect{
				card.Stun{Target: card.Target.FriendlyCreature.ExceptHouse(card.House.Self)},
				card.Stun{Target: card.Target.EnemyCreature.ExceptHouse(card.House.Self)},
			},
		}),
)
