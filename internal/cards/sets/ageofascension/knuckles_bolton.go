package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Knuckles Bolton
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive, Skirmish.
var KnucklesBolton = card.New(
	"Knuckles Bolton",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 271),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive, card.Keyword.Skirmish),
)
