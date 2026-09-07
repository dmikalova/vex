//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GarciasBlaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Variant
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may deal 2D to a creature, or attach Garcia's Blaster to Sensor Chief Garcia."
//	After you attach Garcia's Blaster to Sensor Chief Garcia, steal 1A.
var GarciasBlaster = card.New(
	"Garcia's Blaster",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 347),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
