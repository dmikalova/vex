//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LightOfTheArchons
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature gets +1 power and +1 armor for each upgrade attached to it.
var LightOfTheArchons = card.New(
	"Light of the Archons",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, 300),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
