//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// MartyrsEnd
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy any number of friendly creatures. Gain 1A for each creature destroyed this way.
var MartyrsEnd = card.New(
	"Martyr's End",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 255),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
