//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ghostform
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gains invulnerable. (It cannot be destroyed or dealt damage.)
//	This creature gains, "Fight/Reap: Archive Ghostform."
var Ghostform = card.New(
	"Ghostform",
	card.House.Brobnar,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A01"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
