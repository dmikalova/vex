package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Daemo-Beast
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Beast
//
//	Skirmish.
//	Destroyed: Steal 1 Æmber.
var DaemoBeast = set.New(
	"Daemo-Beast",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "363"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Beast),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Destroyed, card.StealAember{Amount: 1}),
)
