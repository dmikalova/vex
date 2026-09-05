package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Mars First
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Ready and use a friendly Mars creature.
var MarsFirst = card.New(
	"Mars First",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 165),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature.OfHouse(card.House.Self),
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.UseVerb{}},
		}),
)
