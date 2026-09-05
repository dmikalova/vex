package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Eye of Judgment
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Action: Purge a creature from a discard pile.
var EyeOfJudgment = card.New(
	"Eye of Judgment",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 253),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.PurgeCard{
			Zone: card.Discard,
			Type: card.Type.Creature,
		}),
)
