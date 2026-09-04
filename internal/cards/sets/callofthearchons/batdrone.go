package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Batdrone
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Robot
//
//	Skirmish.
//	Fight: Steal 1 Æmber.
var Batdrone = card.New(
	"Batdrone",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 136),
	card.WithPower(2),
	card.WithTraits(card.Traits.Robot),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Fight, card.StealAember{Amount: 1}),
)
