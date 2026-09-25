package massmutation

import "github.com/dmikalova/vex/internal/card"

// Scaethe
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Destroyed: Destroy the least powerful enemy creature.
var Scaethe = set.New(
	"Scaethe",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "052"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Destroyed, card.Destroy{
			Target: card.Target.EachEnemyCreature.Refine(card.LeastPowerful),
		}),
)
