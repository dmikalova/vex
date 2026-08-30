//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Shaffles
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	At the end of your turn, your opponent loses 1 Aember.
var Shaffles = card.New(
	"Shaffles",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 95),
	card.WithPower(2),
	card.WithTraits("Imp"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
