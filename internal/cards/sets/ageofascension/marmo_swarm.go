//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// MarmoSwarm
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast
//
//	Marmo Swarm gets +1 power for each A in your pool.
var MarmoSwarm = card.New(
	"Marmo Swarm",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 327),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
