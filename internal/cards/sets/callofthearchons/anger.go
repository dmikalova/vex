package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Anger
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Ready and fight with a friendly creature.
var Anger = set.New(
	"Anger",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "1"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature,
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		}),
)
