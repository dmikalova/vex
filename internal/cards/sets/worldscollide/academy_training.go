//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// AcademyTraining
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Rare
//
//	If you control this creature, it belongs to house Logos. (Instead of its original house.)
//	This creature gains, "Reap: Draw a card."
var AcademyTraining = card.New(
	"Academy Training",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 161),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
