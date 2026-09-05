package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Shard of Strength
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, give a friendly creature 3 +1 power counters.
var ShardOfStrength = card.New(
	"Shard of Strength",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 48),
	card.WithTraits(card.Traits.Item, card.Traits.Shard),
	card.WithAbility(
		card.Trigger.Action, card.AddPowerCounter{
			Target: card.Target.FriendlyCreature,
			Amount: 3,
			Per: card.InPlay{
				Player: card.Controller,
				Trait:  card.Traits.Shard,
			},
		}),
)
