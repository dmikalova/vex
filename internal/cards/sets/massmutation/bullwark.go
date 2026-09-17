package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Bull-wark
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Mutant • Knight
//
//	Assault 2.
//	Each neighboring creature gains assault 2.
var Bullwark = set.New(
	"Bull-wark",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "127"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Mutant, card.Traits.Knight),
	card.WithAssault(2),
	card.WithConstant(card.ConstantAbility{
		Target:       card.Target.EachCreature.Neighboring(),
		AssaultBonus: 2,
	}),
)
