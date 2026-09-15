package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Teleporter Chief Tink
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Alien
//
//	Elusive.
//	Action: Swap this creature with another friendly creature in your battleline. Use the other creature.
var TeleporterChiefTink = set.New(
	"Teleporter Chief Tink",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "317"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Alien),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Swap{With: card.Target.OtherFriendlyCreature},
			card.OnChooseCreature{
				Target: card.Target.TheOtherCreature,
				Verbs:  []card.CreatureVerb{card.UseVerb{}},
			},
		}}),
)
