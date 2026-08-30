package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Anger
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Ready and fight with a friendly creature.
var Anger = card.New(
	"Anger",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 1),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature,
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		}),
)
