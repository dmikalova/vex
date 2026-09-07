//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CauldronBoil
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal damage to each creature equal to the amount of damage on that creature.
var CauldronBoil = card.New(
	"Cauldron Boil",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 354),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
