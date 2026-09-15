package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Mindworm
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Beast
//
//	Elusive.
//	Before Fight: Deal damage equal to its power to each neighbor of the creature Mindworm fights.
var Mindworm = set.New(
	"Mindworm",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "168"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.BeforeFight, card.DealDamage{
			AmountFrom: card.PowerOfChosen{},
			Target:     card.Target.CreatureFought.NeighborsOf(),
		}),
)
