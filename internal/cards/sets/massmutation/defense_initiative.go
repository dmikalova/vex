//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DefenseInitiative
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ward a creature. You may exalt that creature. If you exalt it, ward each of its neighbors.
var DefenseInitiative = set.New(
	"Defense Initiative",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "191"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
