//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Painmail
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "After any player chooses Dis as their active house, put Painmail into its owner's archives and destroy this creature."
var Painmail = set.New(
	"Painmail",
	card.House.Dis,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.MM, "042"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
