package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Niffle Grounds
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	Action: Choose a Creature - for the remainder of the turn, it loses taunt and elusive.
var NiffleGrounds = card.New(
	"Niffle Grounds",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "346"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.LoseKeywords{
				Target:   card.Target.Triggering,
				Keywords: []card.KeywordValue{card.Keyword.Taunt, card.Keyword.Elusive},
			},
		}),
)
