package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Hallowed Blaster
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Heal 3 damage from a creature.
var HallowedBlaster = set.New(
	"Hallowed Blaster",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, "233"),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.Heal{
			Amount: 3,
			Target: card.Target.Creature,
		}),
)
