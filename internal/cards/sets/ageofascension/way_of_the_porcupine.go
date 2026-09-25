package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Way of the Porcupine
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains +3 hazardous.
var WayOfThePorcupine = set.New(
	"Way of the Porcupine",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "350"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{HazardousBonus: 3}),
)
