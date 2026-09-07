//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Snaglet
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Action: Choose a house. If your opponent chooses that house as their active house on their next turn, steal 2A.
var Snaglet = card.New(
	"Snaglet",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 118),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
