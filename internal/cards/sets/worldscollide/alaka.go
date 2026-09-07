//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Alaka
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Giant
//
//	If you have used a creature to fight this turn, Alaka enters play ready.
var Alaka = card.New(
	"Alaka",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 1),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
