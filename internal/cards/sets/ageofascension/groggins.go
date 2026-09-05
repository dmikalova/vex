package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Groggins
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  8
//	Traits: Giant
//
//	Groggins can only fight flank creatures.
var Groggins = card.New(
	"Groggins",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 11),
	card.WithPower(8),
	card.WithTraits(card.Traits.Giant),
	card.WithFightRestriction(card.Target.EachCreature.OnFlank()),
)
