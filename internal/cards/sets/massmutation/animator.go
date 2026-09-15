//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Animator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Move an artifact to a flank of its controller's battleline. For the remainder of the turn, it is a creature with 3 power that belongs to the active house. (It leaves the battleline when it's no longer a creature.)
var Animator = set.New(
	"Animator",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "100"),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
