package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Barrister Joya
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  1
//	Traits: Human • Knight
//
//	Enemy Creatures cannot reap.
var BarristerJoya = card.New(
	"Barrister Joya",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "228"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithRestrictions(card.Restrictions{Reaping: card.Opponent}),
)
