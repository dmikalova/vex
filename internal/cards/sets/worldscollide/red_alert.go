//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// RedAlert
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Common
//
//	Play: If there are more enemy creatures than friendly creatures, deal damage to each enemy creature equal to the difference.
var RedAlert = card.New(
	"Red Alert",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 303),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
