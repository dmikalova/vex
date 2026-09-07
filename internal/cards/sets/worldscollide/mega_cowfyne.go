//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MegaCowfyne
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  7
//	Traits: Giant
//
//	Before Fight: Deal 2D to each neighbor of the creature Mega Cowfyne fights.
var MegaCowfyne = card.New(
	"Mega Cowfyne",
	card.House.Brobnar,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "55"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
