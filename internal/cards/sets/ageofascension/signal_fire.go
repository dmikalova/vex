package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Signal Fire
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Signal Fire. For the remainder of the turn, each friendly Brobnar creature may fight.
var SignalFire = card.New(
	"Signal Fire",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 49),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.GrantFightForFriendlyHouse{House: card.House.Self},
		}}),
)
