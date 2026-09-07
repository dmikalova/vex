//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Stunner
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may stun a creature."
var Stunner = card.New(
	"Stunner",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 319),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
