package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Inspiration
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ready and use a friendly creature.
var Inspiration = card.New(
	"Inspiration",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 220),
	card.WithAbility(
		card.Trigger.Play, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature,
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.UseVerb{}},
		}),
)
