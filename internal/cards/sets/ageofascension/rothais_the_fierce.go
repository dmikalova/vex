package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Rothais the Fierce
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Taunt, Hazardous 4.
var RothaisTheFierce = card.New(
	"Rothais the Fierce",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 242),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithHazardous(4),
)
