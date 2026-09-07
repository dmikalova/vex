//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Xenotraining
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each house represented among friendly creatures, a friendly creature captures 1A.
var Xenotraining = card.New(
	"Xenotraining",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "323"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
