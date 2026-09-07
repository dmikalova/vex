//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BarnRazing
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, your opponent loses 1A each time a friendly creature fights.
var BarnRazing = card.New(
	"Barn Razing",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "4"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
