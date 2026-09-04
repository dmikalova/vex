package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Replicator
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Mutant
//
//	Reap: Trigger the reap effect of another creature.
var Replicator = card.New(
	"Replicator",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 150),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Reap, card.TriggerAbility{
			Trigger: card.Trigger.Reap,
			Target:  card.Target.Creature.Other(),
		}),
)
