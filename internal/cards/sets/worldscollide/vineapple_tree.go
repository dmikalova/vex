package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Vineapple Tree
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	Each player's keys cost +1 Æmber for each growth counter on Vineapple Tree.
//	After a player forges a key, remove each growth counter from Vineapple Tree.
//	Action: Put a growth counter on Vineapple Tree.
var VineappleTree = card.New(
	"Vineapple Tree",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "402"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	card.WithKeyCost(
		card.KeyCostChange(card.EachPlayer, 1).Per(card.CountersOnThis{Kind: card.Counter.Growth})),
	card.WithAbility(
		card.Trigger.AfterPlayerForgesKey, card.RemoveCounters{
			Kind:   card.Counter.Growth,
			Target: card.Target.This,
		}),
	card.WithAbility(
		card.Trigger.Action, card.PlaceCounter{
			Kind:   card.Counter.Growth,
			Target: card.Target.This,
		}),
)
