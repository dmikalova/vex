//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Narp
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  8
//	Armor:  1
//	Traits: Giant
//
//	Narp's neighbors cannot reap.
var Narp = card.New(
	"Narp",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 12),
	card.WithPower(8),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
