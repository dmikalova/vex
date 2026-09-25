package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Camouflage
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Creatures not on a flank cannot fight this creature.
var Camouflage = set.New(
	"Camouflage",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "337"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{ProtectsFromNonFlank: true}),
)
