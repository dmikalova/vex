//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// WailOfTheDamned
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Enhance Capture.
//	Play: Destroy a creature with no bonus icons.
var WailOfTheDamned = set.New(
	"Wail of the Damned",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "033"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
