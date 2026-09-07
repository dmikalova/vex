//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MegaAlaka
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Giant
//
//	If you have used a creature to fight this turn, Mega Alaka enters play ready.
var MegaAlaka = card.New(
	"Mega Alaka",
	card.House.Brobnar,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "54"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
