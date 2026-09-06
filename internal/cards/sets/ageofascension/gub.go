package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Gub
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Demon
//
//	Gub gains +5 power and taunt while it is not on a flank.
var Gub = card.New(
	"Gub",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 60),
	card.WithPower(1),
	card.WithTraits(card.Traits.Demon),
	card.WithConstant(card.ConstantAbility{
		Target:        card.Target.This,
		PowerBonus:    5,
		Keywords:      card.Keywords(card.Keyword.Taunt),
		WhileOffFlank: true,
	}),
)
