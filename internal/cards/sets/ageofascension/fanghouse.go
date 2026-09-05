package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Fanghouse
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Assault 3, Hazardous 3.
var Fanghouse = card.New(
	"Fanghouse",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 321),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithAssault(3),
	card.WithHazardous(3),
)
