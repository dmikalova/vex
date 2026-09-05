package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Knoxx
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Knoxx gains +3 power for each neighbor it has.
var Knoxx = card.New(
	"Knoxx",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 326),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This,
		PowerBonus: 3,
		Per:        card.NeighborsOfThis{},
	}),
)
