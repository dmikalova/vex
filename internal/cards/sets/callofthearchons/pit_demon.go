package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pit Demon
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Action: Steal 1 Æmber.
var PitDemon = card.New(
	"Pit Demon",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 92),
	card.WithPower(5),
	card.WithTraits("Demon"),
	card.WithAbility(
		card.Trigger.Action, card.StealAember{Amount: 1}),
)
