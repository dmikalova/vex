//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Toad
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Beast
//
//	Toad cannot reap.
var Toad = card.New(
	"Toad",
	card.House.Untamed,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "405"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
