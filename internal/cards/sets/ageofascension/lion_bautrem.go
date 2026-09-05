package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// "Lion" Bautrem
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Human • Knight
//
//	Deploy.
//	Each neighboring creature gains +2 power.
var LionBautrem = card.New(
	"\"Lion\" Bautrem",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 211),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithConstant(card.ConstantAbility{
		PowerBonus: 2,
		Target:     card.Target.EachCreature.Neighboring(),
	}),
)
