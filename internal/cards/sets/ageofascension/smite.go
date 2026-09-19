package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Smite
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ready and fight with a friendly creature. Deal 2 damage to each neighbor of the fought creature.
var Smite = set.New(
	"Smite",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "224"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.OnChooseCreature{
				Target: card.Target.FriendlyCreature,
				Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
			},
			card.DealDamage{
				Amount: 2,
				Target: card.Target.TheFoughtCreature.NeighborsOf(),
			},
		}}),
)
