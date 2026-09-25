package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Panpaca, Anga
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Beast
//
//	Each creature to the right of Panpaca, Anga gains +2 power.
var PanpacaAnga = set.New(
	"Panpaca, Anga",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "347"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.EachCreature.ToRightOfSource(),
		PowerBonus: 2,
	}),
)
