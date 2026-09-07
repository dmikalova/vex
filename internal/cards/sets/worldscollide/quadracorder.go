//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Quadracorder
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	Your opponent's keys cost +1A for each house represented among friendly creatures (to a maximum of 3).
var Quadracorder = card.New(
	"Quadracorder",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 316),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
