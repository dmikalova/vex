//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GalacticCensus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If there are exactly 3 or exactly 4 houses represented among creatures in play, gain 1A. If there are exactly 5, gain 2A. If there are 6 or more, gain 3A.
var GalacticCensus = card.New(
	"Galactic Census",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "332"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
