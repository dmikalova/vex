package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Bloodshard Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	After a creature reaps, destroy it.
var BloodshardImp = card.New(
	"Bloodshard Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 70),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Destroy{Target: card.Target.Triggering}),
)
