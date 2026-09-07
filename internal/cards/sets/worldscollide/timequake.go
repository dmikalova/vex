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
//	Rarity: FIXED
//	Æmber:  1
//
//	Play: Shuffle each friendly card in play into your deck. Draw a card for each card shuffled into your deck this way.
var Timequake = card.New(
	"Timequake",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 0),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
