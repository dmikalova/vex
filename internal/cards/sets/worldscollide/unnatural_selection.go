//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// UnnaturalSelection
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose 3 friendly creatures and 3 enemy creatures. Destroy each other creature.
var UnnaturalSelection = card.New(
	"Unnatural Selection",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 367),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
