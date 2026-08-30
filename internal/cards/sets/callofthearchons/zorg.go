//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Zorg
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Beast
//
//	Zorg enters play stunned.
//	Before Fight: Stun the creature Zorg fights and each of that creature's neighbors.
var Zorg = card.New(
	"Zorg",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 206),
	card.WithPower(7),
	card.WithTraits("Beast"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
