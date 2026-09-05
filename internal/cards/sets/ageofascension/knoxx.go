//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Knoxx
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Knoxx gets +3 power for each neighbor it has.
var Knoxx = card.New(
	"Knoxx",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 326),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
