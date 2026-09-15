//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// OppositionResearch
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Enhance Damage.
//	Play: Enemy creatures cannot reap during your opponent's next turn.
var OppositionResearch = set.New(
	"Opposition Research",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "077"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
