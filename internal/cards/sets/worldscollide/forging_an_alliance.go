//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ForgingAnAlliance
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Forge a key at +7A current cost, reduced by 1A (to a maximum of 6) for each house represented among cards in play.
var ForgingAnAlliance = card.New(
	"Forging an Alliance",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 331),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
