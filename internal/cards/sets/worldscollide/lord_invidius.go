//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LordInvidius
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Demon • Leader
//
//	Elusive.
//	While Lord Invidius is in the center of your battleline, it gains, "Reap: Take control of an enemy flank creature and exhaust it. While under your control, it belongs to house Dis."
var LordInvidius = card.New(
	"Lord Invidius",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "110"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon, card.Traits.Leader),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
