package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Quiet Anvil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Each player's keys cost -2 Æmber.
//	After a player forges a key, destroy The Quiet Anvil.
var TheQuietAnvil = card.New(
	"The Quiet Anvil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 282),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithKeyCost(card.KeyCostChange(card.EachPlayer, -2)),
	card.WithAbility(
		card.Trigger.AfterPlayerForgesKey, card.Destroy{Target: card.Target.This}),
)
