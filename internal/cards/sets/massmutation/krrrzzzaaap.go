//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Krrrzzzaaap
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy each non-Mutant creature. Gain 1 chain.
var Krrrzzzaaap = set.New(
	"Krrrzzzaaap!!!",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "090"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
