package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Killzord Mk. 9001
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains +2 power, +2 armor, and skirmish.
var KillzordMk9001 = card.New(
	"Killzord Mk. 9001",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 201),
	card.WithStatic(card.StaticModifier{
		PowerBonus: 2,
		ArmorBonus: 2,
		Keywords:   card.Keywords(card.Keyword.Skirmish),
	}),
)
