package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dodger
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Elf • Thief
//
//	Fight: Steal 1 Æmber.
var Dodger = card.New(
	"Dodger",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 308),
	card.WithPower(5),
	card.WithTraits("Elf", "Thief"),
	card.WithAbility(
		card.Trigger.Fight, card.StealAember{Amount: 1}),
)
