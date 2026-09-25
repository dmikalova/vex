package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Gron Nine-Toes
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Gron Nine-Toes gains +4 power while it is damaged.
var GronNineToes = set.New(
	"Gron Nine-Toes",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "9"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This.Damaged(),
		PowerBonus: 4,
	}),
)
