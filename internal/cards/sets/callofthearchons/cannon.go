package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Cannon
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Action: Deal 2 damage to a creature.
var Cannon = card.New(
	"Cannon",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 21),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.DealDamage{
			Amount: 2,
			Target: card.Target.Creature,
		}),
)
