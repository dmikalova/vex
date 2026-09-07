package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Special Agent "Fingers"
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Steal 1 Æmber.
var SpecialAgentFingers = card.New(
	"Special Agent \"Fingers\"",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "339"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.StealAember{Amount: 1}),
)
