package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Seraphic Armor
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains +1 armor.
//	Play: Fully heal this creature.
var SeraphicArmor = card.New(
	"Seraphic Armor",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 263),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{ArmorBonus: 1}),
	card.WithAbility(
		card.Trigger.Play, card.Heal{Fully: true, Target: card.Target.This}),
)
