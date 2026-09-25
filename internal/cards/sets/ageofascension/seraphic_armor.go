package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Seraphic Armor
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains +1 armor.
//	Play: Fully heal this creature.
var SeraphicArmor = set.New(
	"Seraphic Armor",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "263"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{ArmorBonus: 1}),
	card.WithAbility(
		card.Trigger.Play, card.Heal{
			Fully:  true,
			Target: card.Target.This,
		}),
)
