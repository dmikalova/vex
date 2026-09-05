package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Containment Field
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "After this creature is used, destroy this creature."
var ContainmentField = card.New(
	"Containment Field",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 178),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.UsedSelf,
			Effect:  card.Destroy{Target: card.Target.This},
		}},
	}),
)
