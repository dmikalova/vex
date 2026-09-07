//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// WallsBlaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Variant
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may deal 2D to a creature, or attach Walls' Blaster to Chief Engineer Walls."
//	After you attach Walls' Blaster to Chief Engineer Walls, stun a creature for each upgrade on Chief Engineer Walls.
var WallsBlaster = card.New(
	"Walls' Blaster",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 352),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
