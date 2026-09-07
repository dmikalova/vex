//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// FyreBreath
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gets +3 power and gains, "Before Fight: Deal 2D to each neighbor of the creature this creature fights."
var FyreBreath = card.New(
	"Fyre-Breath",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 20),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
