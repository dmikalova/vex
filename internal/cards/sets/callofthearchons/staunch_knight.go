package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Staunch Knight
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Staunch Knight gains +2 power while it is on a flank.
var StaunchKnight = card.New(
	"Staunch Knight",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 259),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits("Human", "Knight"),
	card.WithConstant(card.ConstantAbility{
		PowerBonus: 2,
		Target:     card.Target.This.OnFlank(),
	}),
)
