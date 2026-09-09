package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Snag's Mirror
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	After a player chooses an active house, their opponent cannot choose the same house as their active house on their next turn.
var SnagsMirror = card.New(
	"Snag's Mirror",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "117"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterAnyPlayerChoosesHouse,
		card.ForbidSameActiveHouseNextTurn{},
	),
)
