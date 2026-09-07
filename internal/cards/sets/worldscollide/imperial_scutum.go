//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ImperialScutum
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature gets +2 armor and gains, "Destroyed: Move each A on this creature to the common supply."
var ImperialScutum = card.New(
	"Imperial Scutum",
	card.House.Saurian,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, 185),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
