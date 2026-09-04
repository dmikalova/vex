package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dominator Bauble
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	Action: Use a friendly creature.
var DominatorBauble = card.New(
	"Dominator Bauble",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, 73),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature,
			Verbs:  []card.CreatureVerb{card.UseVerb{}},
		}),
)
