//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Borrow
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Take control of an enemy artifact. While under your control, it belongs to house Shadows. (Instead of its original house.)
var Borrow = set.New(
	"\"Borrow\"",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "263"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
