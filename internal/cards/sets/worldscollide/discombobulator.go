//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Discombobulator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Your A cannot be stolen."
var Discombobulator = card.New(
	"Discombobulator",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 149),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
