package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bulwark
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Each neighboring creature gains +2 armor.
var Bulwark = card.New(
	"Bulwark",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 238),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits("Human", "Knight"),
	card.WithConstant(card.ConstantAbility{
		ArmorBonus: 2,
		Target:     card.Target.EachCreature.Neighboring(),
	}),
)
