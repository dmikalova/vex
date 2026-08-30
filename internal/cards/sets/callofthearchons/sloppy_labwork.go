//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SloppyLabwork
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Archive a card. Discard a card.
var SloppyLabwork = card.New(
	"Sloppy Labwork",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 123),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
