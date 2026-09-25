package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Scientifical Hack
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Equation
//
//	Versatile.
//	Action: Destroy Scientifical Hack. For the remainder of the turn, you may use friendly artifacts as if they belonged to the active house.
var ScientificalHack = set.New(
	"Scientifical Hack",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "154"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Equation),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.MayPlayOrUse{
				Houses: card.GrantHouses.Any,
				Grant:  card.GrantUse,
				Types:  card.Types.Of(card.Type.Artifact),
			},
		}}),
)
