package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Standardized Testing
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each Creature with the lowest power and each Creature with the highest power.
var StandardizedTesting = set.New(
	"Standardized Testing",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "119"),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.Refine(
				card.AnyOf(card.LowestPower, card.HighestPower),
			),
		}),
)
