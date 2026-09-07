//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// AxiomOfGrisk
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ward a creature. Destroy each creature with no A on it. Gain 2 chains.
var AxiomOfGrisk = card.New(
	"Axiom of Grisk",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 182),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
