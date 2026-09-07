//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Valoocanth
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Aquan
//
//	While the tide is low, Valoocanth cannot be used. (By default, the tide is neither high or low.)
//	Fight/Reap: Exhaust an enemy creature and each of its neighbors.
var Valoocanth = card.New(
	"Valoocanth",
	card.House.Brobnar,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A10"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Aquan),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
