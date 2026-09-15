package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Daemo-Knight
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Mutant • Knight
//
//	Destroyed: Steal 1 Æmber.
var DaemoKnight = set.New(
	"Daemo-Knight",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "132"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Destroyed, card.StealAember{Amount: 1}),
)
