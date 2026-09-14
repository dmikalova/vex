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
//	Each neighboring Creature gains +2 armor.
var Bulwark = set.New(
	"Bulwark",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, "238"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithConstant(card.ConstantAbility{
		ArmorBonus: 2,
		Target:     card.Target.EachCreature.Neighboring(),
	}),
)
