//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ghosthawk
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
//	Deploy. (This creature can enter play anywhere in your battleline.)
//	Play: You may reap with each neighboring creature, one at a time.
var Ghosthawk = card.New(
	"Ghosthawk",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 356),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
