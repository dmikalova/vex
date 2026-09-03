package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mooncurser
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Elf • Thief
//
//	Skirmish, Poison.
//	Fight: Steal 1 Æmber.
var Mooncurser = card.New(
	"Mooncurser",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 304),
	card.WithPower(1),
	card.WithTraits("Elf", "Thief"),
	card.WithKeywords(card.Keyword.Skirmish, card.Keyword.Poison),
	card.WithAbility(
		card.Trigger.Fight, card.StealAember{Amount: 1}),
)
