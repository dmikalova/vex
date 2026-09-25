package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Guard Disguise
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Destroy Guard Disguise. If your opponent has 3 Æmber or fewer, steal 3 Æmber.
var GuardDisguise = set.New(
	"Guard Disguise",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "302"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.Conditional{
				Cond: card.PoolAember{
					Player: card.Opponent,
					Is:     card.AtMost,
					Amount: 3,
				},
				Then: card.StealAember{Amount: 3},
			},
		}}),
)
