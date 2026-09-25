package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Hideaway Hole
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	Versatile.
//	Action: Destroy Hideaway Hole. Each friendly creature gains elusive until the start of your next turn.
var HideawayHole = set.New(
	"Hideaway Hole",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "287"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.GainKeywords{
				Target:   card.Target.EachFriendlyCreature,
				Keywords: []card.KeywordValue{card.Keyword.Elusive},
				Duration: card.Duration.StartOfPlayerNextTurn,
			},
		}}),
)
