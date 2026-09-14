package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Containment Field
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This Creature gains, "After this Creature is used, destroy this Creature."
var ContainmentField = set.New(
	"Containment Field",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "178"),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.UsedSelf,
			Effect:  card.Destroy{Target: card.Target.This},
		}},
	}),
)
