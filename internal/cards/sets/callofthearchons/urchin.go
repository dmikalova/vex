package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Urchin
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Play: Steal 1 Æmber.
var Urchin = card.New(
	"Urchin",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 315),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{Amount: 1}),
)
