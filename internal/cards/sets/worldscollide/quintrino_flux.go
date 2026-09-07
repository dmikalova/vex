//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// QuintrinoFlux
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a friendly creature and an enemy creature. Destroy the chosen creatures and each creature with the same power as either of the chosen creatures.
var QuintrinoFlux = card.New(
	"Quintrino Flux",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 317),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
