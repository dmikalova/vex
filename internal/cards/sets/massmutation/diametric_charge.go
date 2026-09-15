//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DiametricCharge
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 1D to a creature, with 2D splash.
var DiametricCharge = set.New(
	"Diametric Charge",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "070"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
