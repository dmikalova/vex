//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DataForge
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: You may forge a key at +10A current cost, reduced by 1A for each card in your hand.
var DataForge = card.New(
	"Data Forge",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 148),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
