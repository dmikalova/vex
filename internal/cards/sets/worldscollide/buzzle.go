//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Buzzle
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Skirmish. (When you use this creature to fight, it is dealt no damage in return.)
//	Play/Fight: You may purge one of Buzzle's neighbors. If you do, ready Buzzle.
var Buzzle = card.New(
	"Buzzle",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 70),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
