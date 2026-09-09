package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Gron Nine-Toes
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  7
//	Traits: Giant
//
//	Mega Gron Nine-Toes gains +4 power while it is damaged.
var MegaGronNineToes = card.New(
	"Mega Gron Nine-Toes",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "58"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This.Damaged(),
		PowerBonus: 4,
	}),
)
