package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Aember Conduction Unit
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	After an enemy Creature reaps, if it is the first time a Creature has reaped this turn, stun it.
var AemberConductionUnit = card.New(
	"Aember Conduction Unit",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "176"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterEnemyCreatureReaps, card.Conditional{
			Cond: card.FirstReapOfTurn{},
			Then: card.Stun{Target: card.Target.Triggering},
		}),
)
