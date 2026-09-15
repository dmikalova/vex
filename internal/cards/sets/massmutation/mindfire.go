//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mindfire
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent discards a random card from their hand. Steal 1A for each bonus icon on the discarded card.
var Mindfire = set.New(
	"Mindfire",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "012"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
