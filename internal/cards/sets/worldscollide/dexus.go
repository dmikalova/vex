//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Dexus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Variant
//	Power:  5
//	Traits: Demon
//
//	After your opponent plays a creature on their right flank, they lose 1A.
var Dexus = card.New(
	"Dexus",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 124),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
