//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MolinasBlaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may deal 2D to a creature, or attach Molina's Blaster to Armsmaster Molina."
//	After you attach Molina's Blaster to Armsmaster Molina, you may deal 3D to a creature.
var MolinasBlaster = card.New(
	"Molina's Blaster",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, 302),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
