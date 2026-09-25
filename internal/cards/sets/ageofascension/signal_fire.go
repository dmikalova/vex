package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Signal Fire
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Signal Fire. For the remainder of the turn, each friendly Brobnar creature may fight.
var SignalFire = set.New(
	"Signal Fire",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "49"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.MayPlayOrUse{
				Houses: card.GrantHouses.Named(card.House.Self),
				Grant:  card.GrantFight,
			},
		}}),
)
