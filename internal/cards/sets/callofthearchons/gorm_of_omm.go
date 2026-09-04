package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gorm of Omm
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Gorm of Omm and an artifact.
var GormOfOmm = card.New(
	"Gorm of Omm",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 232),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.Destroy{Target: card.Target.Artifact},
		}}),
)
