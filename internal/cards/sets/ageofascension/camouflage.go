package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Camouflage
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	Creatures not on a flank cannot fight this creature.
var Camouflage = card.New(
	"Camouflage",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "337"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{ProtectsFromNonFlank: true}),
)
