package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Universal Keylock
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Each player's keys cost +3 Æmber.
//	After a player forges a key, destroy Universal Keylock.
var UniversalKeylock = set.New(
	"Universal Keylock",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "178"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithKeyCost(card.KeyCostChange(card.EachPlayer, 3)),
	card.WithAbility(
		card.Trigger.AfterPlayerForgesKey, card.Destroy{Target: card.Target.This}),
)
