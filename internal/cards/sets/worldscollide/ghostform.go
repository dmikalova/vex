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
//	Rarity: FIXED
//	Æmber:  1
//
//	This creature gains invulnerable. (It cannot be destroyed or dealt damage.)
//	This creature gains, "Fight/Reap: Archive Ghostform."
var Ghostform = card.New(
	"Ghostform",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 0),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
