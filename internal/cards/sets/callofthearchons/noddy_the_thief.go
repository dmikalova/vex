package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Noddy the Thief
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Steal 1 Æmber.
var NoddyTheThief = card.New(
	"Noddy the Thief",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 306),
	card.WithPower(2),
	card.WithTraits("Elf", "Thief"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.StealAember{Amount: 1}),
)
