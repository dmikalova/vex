//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ShoulderArmor
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	While this creature is on a flank, it gets +2 armor and +2 power.
var ShoulderArmor = card.New(
	"Shoulder Armor",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 266),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
