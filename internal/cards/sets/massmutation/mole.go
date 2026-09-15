//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mole
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Your opponent may spend A on this creature as if it were in their pool."
var Mole = set.New(
	"Mole",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.MM, "287"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
