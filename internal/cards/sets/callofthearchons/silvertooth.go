package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Silvertooth
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Silvertooth enters play ready.
var Silvertooth = card.New(
	"Silvertooth",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 311),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithEntersPlay(card.Ready{Target: card.Target.This}),
)
