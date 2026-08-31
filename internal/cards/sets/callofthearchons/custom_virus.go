package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Custom Virus
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Custom Virus. Purge a creature from your hand. Destroy each creature that shares a trait with the purged creature.
var CustomVirus = card.New(
	"Custom Virus",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 183),
	card.WithAemberBonus(1),
	card.WithTraits("Weapon"),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Sentence{Effect: card.Destroy{Target: card.Target.This}},
			card.PurgeHandThenDestroyShared{},
		}}),
)
