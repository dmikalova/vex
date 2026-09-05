//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// GoldenAura
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a creature. Fully heal the chosen creature. For the remainder of the turn, the chosen creature is considered to be in house Sanctum and cannot be dealt damage.
var GoldenAura = card.New(
	"Golden Aura",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 217),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
