package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Hideaway Hole
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	Versatile.
//	Action: Destroy Hideaway Hole. Each friendly Creature gains elusive until the start of your next turn.
var HideawayHole = set.New(
	"Hideaway Hole",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "287"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.GainKeyword{
				Target:  card.Target.EachFriendlyCreature,
				Keyword: card.Keyword.Elusive,
			},
		}}),
)
