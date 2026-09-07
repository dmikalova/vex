//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// PlasmaNozzle
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Before Fight: Deal 2D to the attacked creature, with 2D splash."
var PlasmaNozzle = card.New(
	"Plasma Nozzle",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 336),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
