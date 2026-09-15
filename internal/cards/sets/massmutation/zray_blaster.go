//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ZRayBlaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gets +3 power and gains, "Before Fight: Deal 3D to each neighbor of the creature this creature fights."
var ZRayBlaster = set.New(
	"Z-Ray Blaster",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Special,
	card.Provenance(card.MM, "355"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
