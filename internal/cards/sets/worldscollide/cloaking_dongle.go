//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CloakingDongle
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature and each of its neighbors gains elusive. (The first time this creature is attacked each turn, no damage is dealt.)
var CloakingDongle = card.New(
	"Cloaking Dongle",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, "294"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
