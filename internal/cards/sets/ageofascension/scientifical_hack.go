package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Scientifical Hack
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Equation
//
//	Versatile.
//	Action: Destroy Scientifical Hack. For the remainder of the turn, you may use friendly artifacts as if they belonged to the active house.
var ScientificalHack = card.New(
	"Scientifical Hack",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "154"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Equation),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.MayPlayOrUse{
				Houses: card.Houses.Any,
				Grant:  card.GrantUse,
				Types:  card.Types.Artifacts,
			},
		}}),
)
