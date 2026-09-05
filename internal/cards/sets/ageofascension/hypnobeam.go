package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Hypnobeam
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Take control of an enemy creature.
var Hypnobeam = card.New(
	"Hypnobeam",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 181),
	card.WithAbility(
		card.Trigger.Play, card.TakeControl{
			Target:   card.Target.EnemyCreature,
			Duration: card.Duration.Forever,
		}),
)
