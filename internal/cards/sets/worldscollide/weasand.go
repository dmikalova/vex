//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Weasand
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Thief
//
//	Deploy. Elusive.
//	If Weasand is on a flank, destroy it.
//	After a player forges a key, gain 2A.
var Weasand = card.New(
	"Weasand",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 285),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
