//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// XanthyxHarvester
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Xanthyx Harvester cannot be used while it has a non-Mars neighbor.
//	Reap: Gain 1A.
var XanthyxHarvester = card.New(
	"Xanthyx Harvester",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 173),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
