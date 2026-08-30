//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// OverlordGreking
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	After an enemy creature is destroyed fighting Overlord Greking, put that creature into play under your control.
var OverlordGreking = card.New(
	"Overlord Greking",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 87),
	card.WithPower(7),
	card.WithTraits("Demon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
