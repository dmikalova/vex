//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Timequake
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Shuffle each friendly card in play into your deck. Draw a card for each card shuffled into your deck this way.
var Timequake = card.New(
	"Timequake",
	card.House.Brobnar,
	card.Type.Tactic,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A09"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
