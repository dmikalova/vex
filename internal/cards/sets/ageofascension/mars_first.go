package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Mars First
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Ready and use a friendly Mars creature.
var MarsFirst = set.New(
	"Mars First",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "165"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature.House(card.Houses.Named(card.House.Self)),
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.UseVerb{}},
		}),
)
