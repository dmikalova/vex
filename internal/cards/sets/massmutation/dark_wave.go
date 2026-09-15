//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DarkWave
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2D to each non-Mutant
//	creature.
var DarkWave = set.New(
	"Dark Wave",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "247"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
