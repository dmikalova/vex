//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BlastShielding
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature gets +2 armor.
//	After this creature is used, its controller may attach Blast Shielding to one of this creature's neighbors.
var BlastShielding = set.New(
	"Blast Shielding",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.MM, "303"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
