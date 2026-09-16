package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Seeker of Truth
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Armor:  1
//	Traits: Human
//
//	Fight: You may use a friendly non-Sanctum creature.
var SeekerOfTruth = set.New(
	"Seeker of Truth",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "138"),
	card.WithPower(3),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(
		card.Trigger.Fight, card.May{
			Do: card.OnChooseCreature{
				Target: card.Target.FriendlyCreature.House(
					card.Houses.Except(card.House.Self),
				),
				Verbs: []card.CreatureVerb{card.UseVerb{}},
			},
		}),
)
