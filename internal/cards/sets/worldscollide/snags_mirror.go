package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Snag's Mirror
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	After a player chooses an active house, their opponent cannot choose the same house as their active house during their next turn.
var SnagsMirror = set.New(
	"Snag's Mirror",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "117"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterAnyPlayerChoosesHouse,
		card.CannotChooseHouse{
			Player:    card.Opponent,
			Reference: card.JustChosenActiveHouse,
		},
	),
)
