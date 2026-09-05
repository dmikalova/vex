package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Commander Remiel
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Knight
//
//	Reap: Use a friendly non-Sanctum creature.
var CommanderRemiel = card.New(
	"Commander Remiel",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 241),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Reap, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature.ExceptHouse(card.House.Self),
			Verbs:  []card.CreatureVerb{card.UseVerb{}},
		}),
)
