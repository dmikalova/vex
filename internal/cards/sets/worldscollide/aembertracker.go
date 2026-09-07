//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Aembertracker
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast
//
//	Play: Deal 2D to each enemy creature with A on it. This damage cannot be prevented by armor.
var Aembertracker = card.New(
	"Aembertracker",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "324"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
