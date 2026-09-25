package worldscollide

import "github.com/dmikalova/vex/internal/card"

// The Quiet Anvil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Each player's keys cost -2 Æmber.
//	After a player forges a key, destroy The Quiet Anvil.
var TheQuietAnvil = set.New(
	"The Quiet Anvil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "282"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithKeyCost(card.KeyCostChange(card.EachPlayer, -2)),
	card.WithAbility(
		card.Trigger.AfterPlayerForgesKey, card.Destroy{Target: card.Target.This}),
)
