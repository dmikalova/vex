//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// EldestBear
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Beast • Leader • Witch
//
//	Assault 3.
//	While Eldest Bear is in the center of your battleline, it gains, "Before Fight: Gain 2A."
var EldestBear = card.New(
	"Eldest Bear",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 388),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast, card.Traits.Leader, card.Traits.Witch),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
