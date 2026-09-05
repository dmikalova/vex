//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// CollectorWorm
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  5
//	Traits: Beast
//
//	Fight: Archive the creature Collector Worm fights. If that creature leaves your archives, put it in its owner's hand instead.
var CollectorWorm = card.New(
	"Collector Worm",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 162),
	card.WithPower(2),
	card.WithArmor(5),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
