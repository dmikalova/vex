package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Way of the Porcupine
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains +3 hazardous.
var WayOfThePorcupine = card.New(
	"Way of the Porcupine",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 350),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{HazardousBonus: 3}),
)
