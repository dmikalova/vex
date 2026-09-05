package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Spangler Box
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Graft a creature from play, and your opponent gains control of Spangler Box.
//	Destroyed: Put each card under Spangler Box into play under its owner's control.
var SpanglerBox = card.New(
	"Spangler Box",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 132),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Graft{Target: card.Target.Creature},
			card.TakeControl{
				Target:     card.Target.This,
				ToOpponent: true,
				Duration:   card.Duration.Forever,
			},
		}}),
	card.WithAbility(
		card.Trigger.Destroyed, card.PutUnderIntoPlay{}),
)
