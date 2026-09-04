package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Key to Dis
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Key to Dis and each creature.
var KeyToDis = card.New(
	"Key to Dis",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 74),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.Destroy{Target: card.Target.EachCreature},
		}}),
)
