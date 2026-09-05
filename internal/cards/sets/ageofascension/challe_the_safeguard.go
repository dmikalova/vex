package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Challe the Safeguard
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Deploy, Taunt.
var ChalleTheSafeguard = card.New(
	"Challe the Safeguard",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 216),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithKeywords(card.Keyword.Deploy, card.Keyword.Taunt),
)
