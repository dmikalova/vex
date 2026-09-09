package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Universal Keylock
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Each player's keys cost +3 Æmber.
//	After a player forges a key, destroy Universal Keylock.
var UniversalKeylock = card.New(
	"Universal Keylock",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "178"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithKeyCost(card.KeyCostChange(card.EachPlayer, 3)),
	card.WithAbility(
		card.Trigger.AfterPlayerForgesKey, card.Destroy{Target: card.Target.This}),
)
