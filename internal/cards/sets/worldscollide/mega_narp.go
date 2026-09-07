//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MegaNarp
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: FIXED
//	Power:  10
//	Armor:  1
//	Traits: Giant
//
//	Mega Narp's neighbors cannot reap.
var MegaNarp = card.New(
	"Mega Narp",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 60),
	card.WithPower(10),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
