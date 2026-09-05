package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Abond the Armorsmith
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Each other friendly creature gains +1 armor.
//	Action: For the remainder of the turn, each other friendly creature gains +1 armor.
var AbondTheArmorsmith = card.New(
	"Abond the Armorsmith",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 212),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithConstant(card.ConstantAbility{
		ArmorBonus: 1,
		Target:     card.Target.EachOtherFriendlyCreature,
	}),
	card.WithAbility(
		card.Trigger.Action, card.GainStats{
			Target: card.Target.EachOtherFriendlyCreature,
			Armor:  1,
		}),
)
