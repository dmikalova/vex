package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gauntlet of Command
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	Action: Ready and fight with a friendly creature.
var GauntletOfCommand = card.New(
	"Gauntlet of Command",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, 22),
	card.WithTraits("Item"),
	card.WithAbility(
		card.Trigger.Action, card.OnChooseCreature{
			Target: card.Target.FriendlyCreature,
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		}),
)
