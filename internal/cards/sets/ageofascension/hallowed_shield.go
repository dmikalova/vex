package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Hallowed Shield
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	Action: For the remainder of the turn, a creature cannot be dealt damage.
var HallowedShield = card.New(
	"Hallowed Shield",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.AoA, 218),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.PreventDamage{
			Target:   card.Target.Creature,
			Duration: card.Duration.EndOfTurn,
		}),
)
