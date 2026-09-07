//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tribute
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: The most powerful friendly creature captures 2A. You may exalt that creature to repeat the preceding effect.
var Tribute = card.New(
	"Tribute",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "196"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
