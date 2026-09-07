//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// KirbysBlaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Variant
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may deal 2D to a creature, or attach Kirby's Blaster to Com. Officer Kirby."
//	After you attach Kirby's Blaster to Com. Officer Kirby, draw 2 cards.
var KirbysBlaster = card.New(
	"Kirby's Blaster",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 350),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
