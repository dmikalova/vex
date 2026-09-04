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
//	Action: Destroy Custom Virus. Purge a creature from your hand. Destroy each creature that shares a trait with it.
var CustomVirus = card.New(
	"Custom Virus",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 183),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Weapon),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.PurgeCreatureFromHand{},
			card.Destroy{Target: card.Target.EachCreature.SharingTrait()},
		}}),
)
