package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ritual of the Hunt
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Ritual of the Hunt. For the remainder of the turn, you may use friendly Untamed creatures.
var RitualOfTheHunt = card.New(
	"Ritual of the Hunt",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 343),
	card.WithTraits("Power"),
	card.WithAemberBonus(1),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(card.Trigger.Action, card.Sentences{Effects: []card.Effect{
		card.Destroy{Target: card.Target.This},
		card.MayUseFriendlyHouse{House: card.House.Self},
	}}),
)
