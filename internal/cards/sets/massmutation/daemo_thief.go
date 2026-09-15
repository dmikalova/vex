package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Daemo-Thief
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant • Thief
//
//	Elusive.
//	Destroyed: Steal 1 Æmber.
var DaemoThief = set.New(
	"Daemo-Thief",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "246"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Destroyed, card.StealAember{Amount: 1}),
)
