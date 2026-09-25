package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Mega Gron Nine-Toes
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Giant
//
//	Mega Gron Nine-Toes gains +4 power while it is damaged.
var MegaGronNineToes = set.New(
	"Mega Gron Nine-Toes",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "58"),
	card.InCluster(card.Pulled(gronsBrewCluster, 1, 1.25)),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This.Damaged(),
		PowerBonus: 4,
	}),
)
