//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// KhrkharsBlaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Variant
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may deal 2D to a creature, or attach Khrkhar's Blaster to Lieutenant Khrkhar."
//	After you attach Khrkhar's Blaster to Lieutenant Khrkhar, ward Lieutenant Khrkhar.
var KhrkharsBlaster = card.New(
	"Khrkhar's Blaster",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 349),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
